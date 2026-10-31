package spec031_human_task_test

import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcknowledgementCheckpointAndResume(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-ack"

	startWaiting(t, flow, env, runID, "acknowledgement.yaml")
	status := waitForStatus(t, flow, env, runID, "acknowledgement.yaml", "Waiting")
	status.ExpectStdoutNotContains("form:")
	require.Contains(t, status.Stdout(), "maintenance_started")
	require.Contains(t, status.Stdout(), "Confirm that maintenance has started")

	result := complete(t, flow, env, runID, "maintenance_started", "acknowledgement.yaml")
	result.ExpectExitCode(0)
	result.ExpectStdout("Completed human task maintenance_started; DAG-run queued for resume.\n")
	result.ExpectStderr("")

	waitForStatus(t, flow, env, runID, "acknowledgement.yaml", "Succeeded")
	waitForFileContent(t, flow.ProjectPath("continued.txt"), "continued\n")

	repeat := complete(t, flow, env, runID, "maintenance_started", "acknowledgement.yaml")
	repeat.ExpectExitCode(0)
	repeat.ExpectStdout("Human task maintenance_started was already completed.\n")
	repeat.ExpectStderr("")
}

// A partial foreach can leave an independent human checkpoint waiting while
// completion resumes the newly eligible branch without reopening that gate.
func TestPartialCheckpointResume(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-partial-foreach"
	const file = "partial_foreach_checkpoint.yaml"

	startWaiting(t, flow, env, runID, file)
	status := waitForStatus(t, flow, env, runID, file, "Waiting")
	require.Contains(t, status.Stdout(), "partial")

	completed := complete(t, flow, env, runID, "review", file)
	completed.ExpectExitCode(0)
	completed.ExpectStdout("Completed human task review; DAG-run queued for resume.\n")
	completed.ExpectStderr("")
	waitForFileContent(t, flow.ProjectPath("resumed.txt"), "resumed\n")
	status = waitForStatus(t, flow, env, runID, file, "Waiting")
	require.Contains(t, status.Stdout(), "other")

	completed = complete(t, flow, env, runID, "other", file)
	completed.ExpectExitCode(0)
	completed.ExpectStdout("Completed human task other; DAG-run queued for resume.\n")
	completed.ExpectStderr("")
	status = waitForStatus(t, flow, env, runID, file, "Result: Partially Succeeded")
	require.Regexp(t, `(?m)other .*\[succeeded\]`, status.Stdout())
	require.Equal(t, "run\n", readFileContent(t, flow.ProjectPath("foreach-a.out")))
	require.Equal(t, "run\n", readFileContent(t, flow.ProjectPath("foreach-b.out")))
	waitForFileContent(t, flow.ProjectPath("finished.txt"), "finished\n")

	retry := flow.RunWithEnv(env, "retry", "--run-id="+runID, fixtureDAGName(file))
	retry.ExpectExitCode(0)
	status = waitForStatus(t, flow, env, runID, file, "Waiting")
	require.Contains(t, status.Stdout(), "review")
	require.Equal(t, "run\nrun\n", readFileContent(t, flow.ProjectPath("foreach-a.out")))
	require.Equal(t, "run\nrun\n", readFileContent(t, flow.ProjectPath("foreach-b.out")))
	require.Regexp(t, `(?m)other .*\[succeeded\]`, status.Stdout())
}

func readFileContent(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path) // #nosec G304 -- the caller supplies a test-project path.
	require.NoError(t, err)
	return string(data)
}

func TestTypedFormDefaultsAndOutputs(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-form"

	startWaiting(t, flow, env, runID, "typed_form.yaml")
	status := waitForStatus(t, flow, env, runID, "typed_form.yaml", "Waiting")
	require.Contains(t, status.Stdout(), "release_review")
	require.Contains(t, status.Stdout(), "Choose the release target")
	normalizedStatus := collapseStatusLines(status.Stdout())
	require.Contains(t, normalizedStatus, `{"type":"object"`)
	require.Contains(t, normalizedStatus, `"additionalProperties":false`)
	require.Contains(t, normalizedStatus, `"title":"Release review"`)
	require.Contains(t, normalizedStatus, `"environment":{"type":"string","title":"Environment","enum":["staging","production"]}`)
	require.Contains(t, normalizedStatus, `"notify":{"type":"boolean","default":true}`)
	require.Contains(t, normalizedStatus, `"replicas":{"type":"integer","default":2,"minimum":1}`)

	result := complete(t, flow, env, runID, "release_review", "typed_form.yaml", "--input=environment=production")
	result.ExpectExitCode(0)
	result.ExpectStdout("Completed human task release_review; DAG-run queued for resume.\n")
	result.ExpectStderr("")

	waitForStatus(t, flow, env, runID, "typed_form.yaml", "Succeeded")
	waitForFileContent(t, flow.ProjectPath("release.txt"), "production,2,true\n")

	repeat := complete(t, flow, env, runID, "release_review", "typed_form.yaml", "--inputs-json", `{"notify":true,"environment":"production","replicas":2}`)
	repeat.ExpectExitCode(0)
	repeat.ExpectStdout("Human task release_review was already completed.\n")
	repeat.ExpectStderr("")
}

func TestMultipleAndSequentialCheckpoints(t *testing.T) {
	t.Run("independent ordinary branch", func(t *testing.T) {
		flow := harness.NewRunner(t)
		env := sharedEnv(t)
		const runID = "spec031-independent"

		startWaiting(t, flow, env, runID, "independent_branch.yaml")
		waitForFileContent(t, flow.ProjectPath("independent.txt"), "independent\n")
		flow.ExpectNoFile("resumed.txt")

		result := complete(t, flow, env, runID, "review", "independent_branch.yaml")
		result.ExpectExitCode(0)
		result.ExpectStdout("Completed human task review; DAG-run queued for resume.\n")
		result.ExpectStderr("")
		waitForStatus(t, flow, env, runID, "independent_branch.yaml", "Succeeded")
		waitForFileContent(t, flow.ProjectPath("resumed.txt"), "resumed\n")
	})

	t.Run("independent tasks", func(t *testing.T) {
		flow := harness.NewRunner(t)
		env := sharedEnv(t)
		const runID = "spec031-multiple"

		startWaiting(t, flow, env, runID, "multiple_waiting.yaml")
		first := complete(t, flow, env, runID, "review_a", "multiple_waiting.yaml")
		first.ExpectExitCode(0)
		first.ExpectStdout("Completed human task review_a; DAG-run remains waiting.\n")
		first.ExpectStderr("")
		flow.ExpectNoFile("deployed.txt")

		second := complete(t, flow, env, runID, "review_b", "multiple_waiting.yaml")
		second.ExpectExitCode(0)
		second.ExpectStdout("Completed human task review_b; DAG-run queued for resume.\n")
		second.ExpectStderr("")
		waitForStatus(t, flow, env, runID, "multiple_waiting.yaml", "Succeeded")
		waitForFileContent(t, flow.ProjectPath("deployed.txt"), "deployed\n")
	})

	t.Run("independent task branches", func(t *testing.T) {
		flow := harness.NewRunner(t)
		env := sharedEnv(t)
		const runID = "spec031-independent-tasks"

		startWaiting(t, flow, env, runID, "independent_waiting.yaml")
		first := complete(t, flow, env, runID, "review_a", "independent_waiting.yaml")
		first.ExpectExitCode(0)
		first.ExpectStdout("Completed human task review_a; DAG-run queued for resume.\n")
		first.ExpectStderr("")
		waitForFileContent(t, flow.ProjectPath("after_a.txt"), "after_a\n")
		// review_b renders as waiting while the resumed attempt runs, so match the run result.
		status := waitForStatus(t, flow, env, runID, "independent_waiting.yaml", "Result: Waiting")
		require.Contains(t, status.Stdout(), "Confirm B in "+runID)
		flow.ExpectNoFile("finished.txt")

		second := complete(t, flow, env, runID, "review_b", "independent_waiting.yaml")
		second.ExpectExitCode(0)
		second.ExpectStdout("Completed human task review_b; DAG-run queued for resume.\n")
		second.ExpectStderr("")
		waitForStatus(t, flow, env, runID, "independent_waiting.yaml", "Succeeded")
		waitForFileContent(t, flow.ProjectPath("finished.txt"), "finished\n")
	})

	t.Run("sequential tasks", func(t *testing.T) {
		flow := harness.NewRunner(t)
		env := sharedEnv(t)
		const runID = "spec031-sequential"

		startWaiting(t, flow, env, runID, "sequential_waiting.yaml")
		first := complete(t, flow, env, runID, "first_review", "sequential_waiting.yaml")
		first.ExpectExitCode(0)
		first.ExpectStdout("Completed human task first_review; DAG-run queued for resume.\n")
		first.ExpectStderr("")
		status := waitForStatus(t, flow, env, runID, "sequential_waiting.yaml", "Waiting")
		require.Contains(t, status.Stdout(), "second_review")
		waitForFileContent(t, flow.ProjectPath("between.txt"), "between\n")

		second := complete(t, flow, env, runID, "second_review", "sequential_waiting.yaml")
		second.ExpectExitCode(0)
		second.ExpectStdout("Completed human task second_review; DAG-run queued for resume.\n")
		second.ExpectStderr("")
		waitForStatus(t, flow, env, runID, "sequential_waiting.yaml", "Succeeded")
		waitForFileContent(t, flow.ProjectPath("finished.txt"), "finished\n")
	})
}

func TestStoredPromptFormAndDAGSnapshot(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-snapshot"

	startWaiting(t, flow, env, runID, "prompt_snapshot.yaml")
	status := waitForStatus(t, flow, env, runID, "prompt_snapshot.yaml", "Waiting")
	require.Contains(t, status.Stdout(), "Review release for production as operator in "+runID)
	require.Contains(t, status.Stdout(), `${REVIEWER}`)
	require.Contains(t, status.Stdout(), `${params.target}`)

	changed, err := os.ReadFile(flow.ProjectPath("prompt_snapshot_changed.yaml")) // #nosec G304 -- fixed test fixture path.
	require.NoError(t, err)
	flow.WriteFile("prompt_snapshot.yaml", string(changed))

	stored := waitForStatus(t, flow, env, runID, "prompt_snapshot.yaml", "Waiting")
	require.Contains(t, stored.Stdout(), "Review release for production as operator in "+runID)
	require.NotContains(t, stored.Stdout(), "This prompt must not replace")
	require.NoError(t, os.Remove(flow.ProjectPath("prompt_snapshot.yaml")))

	result := complete(t, flow, env, runID, "review", "prompt_snapshot.yaml", "--input=environment=production")
	result.ExpectExitCode(0)
	result.ExpectStdout("Completed human task review; DAG-run queued for resume.\n")
	result.ExpectStderr("")
	waitForStatus(t, flow, env, runID, "spec031_prompt_snapshot", "Succeeded")
	waitForFileContent(t, flow.ProjectPath("snapshot.txt"), "production\n")
}

func TestAdditionalPropertiesDoNotCreateOutputs(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-additional"

	startWaiting(t, flow, env, runID, "additional_properties.yaml")
	input := `{"metadata":{"team":"platform","flags":[true,null]},"environment":"staging"}`
	result := complete(t, flow, env, runID, "collect", "additional_properties.yaml", "--inputs-json", input)
	result.ExpectExitCode(0)
	result.ExpectStdout("Completed human task collect; DAG-run queued for resume.\n")
	result.ExpectStderr("")
	waitForStatus(t, flow, env, runID, "additional_properties.yaml", "Succeeded")
	waitForFileContent(t, flow.ProjectPath("additional.txt"), "staging\n")

	repeat := complete(t, flow, env, runID, "collect", "additional_properties.yaml", "--inputs-json", `{"environment":"staging","metadata":{"flags":[true,null],"team":"platform"}}`)
	repeat.ExpectExitCode(0)
	repeat.ExpectStdout("Human task collect was already completed.\n")
	repeat.ExpectStderr("")
}

func TestConcurrentCompletion(t *testing.T) {
	t.Run("same canonical input", func(t *testing.T) {
		flow := harness.NewRunner(t)
		env := sharedEnv(t)
		const runID = "spec031-concurrent-same"
		scheduler := startWaitingWithScheduler(t, flow, env, runID, "concurrent.yaml")
		defer scheduler.Stop()

		results := runConcurrentCompletions(t, flow, env, runID, "blue", "blue")
		for _, result := range results {
			assert.Equal(t, 0, result.ExitCode(), "stdout:\n%s\nstderr:\n%s", result.Stdout(), result.Stderr())
			assert.Empty(t, result.Stderr())
		}
		waitForStatus(t, flow, env, runID, "concurrent.yaml", "Succeeded")
		waitForFileContent(t, flow.ProjectPath("concurrent.txt"), "blue\n")
	})

	t.Run("different canonical input", func(t *testing.T) {
		flow := harness.NewRunner(t)
		env := sharedEnv(t)
		const runID = "spec031-concurrent-different"
		scheduler := startWaitingWithScheduler(t, flow, env, runID, "concurrent.yaml")
		defer scheduler.Stop()

		results := runConcurrentCompletions(t, flow, env, runID, "blue", "green")
		successes := 0
		failures := 0
		for _, result := range results {
			if result.ExitCode() == 0 {
				successes++
				assert.Empty(t, result.Stderr())
				continue
			}
			failures++
			assert.Empty(t, result.Stdout())
			assert.Contains(t, result.Stderr(), "different input")
		}
		assert.Equal(t, 1, successes)
		assert.Equal(t, 1, failures)
		waitForStatus(t, flow, env, runID, "concurrent.yaml", "Succeeded")
	})
}

func TestPreconditionSkipAndDryRun(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	start := flow.RunWithEnv(env, "start", "--run-id=spec031-skipped", "precondition_skipped.yaml")
	start.ExpectExitCode(0)
	waitForStatus(t, flow, env, "spec031-skipped", "precondition_skipped.yaml", "Succeeded")
	flow.ExpectNoFile("precondition.txt")
	rejected := complete(t, flow, env, "spec031-skipped", "skipped_review", "precondition_skipped.yaml")
	rejected.ExpectNonZeroExitCode()
	rejected.ExpectStdout("")
	rejected.ExpectStderrContains("spec031-skipped", "not waiting")
	rejected.ExpectStderrNotContains("Usage:")

	dry := flow.RunWithEnv(env, "dry", "prompt_snapshot.yaml")
	dry.ExpectExitCode(0)
	flow.ExpectNoFile("snapshot.txt")
}

// A push-back re-runs the rewind target with the feedback, reopens the task
// at the next iteration, rejects a request made for the earlier review, and
// lets completion finish the run once.
func TestPushBackReviewLoop(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-review"
	const file = "push_back.yaml"

	startWaiting(t, flow, env, runID, file)
	waitForFileContent(t, flow.ProjectPath("attempts.txt"), ":\n")
	status := waitForStatus(t, flow, env, runID, file, "Waiting")
	require.Contains(t, status.Stdout(), "push back: rewind to implement")
	require.Contains(t, collapseStatusLines(status.Stdout()), `"feedback":{"type":"string"}`)
	status.ExpectStdoutNotContains("push-back iteration")

	result := pushBack(t, flow, env, runID, "review", file, "--input=feedback=add-tests", "--expected-iteration=0")
	result.ExpectExitCode(0)
	result.ExpectStdout("Pushed back human task review to implement; DAG-run queued for resume.\n")
	result.ExpectStderr("")

	waitForFileContent(t, flow.ProjectPath("attempts.txt"), ":\n1:add-tests\n")
	// The resumed attempt wrote the second line, so the next waiting checkpoint
	// is the reopened review. The review line can show the iteration while that
	// attempt is still finalizing, so wait for the run itself to be waiting.
	status = waitForStatus(t, flow, env, runID, file, "Waiting")
	require.Contains(t, status.Stdout(), "push-back iteration: 1")
	flow.ExpectNoFile("published.txt")

	stale := pushBack(t, flow, env, runID, "review", file, "--input=feedback=again", "--expected-iteration=0")
	stale.ExpectNonZeroExitCode()
	stale.ExpectStdout("")
	stale.ExpectStderrContains("review", "iteration 1")

	done := complete(t, flow, env, runID, "review", file)
	done.ExpectExitCode(0)
	done.ExpectStdout("Completed human task review; DAG-run queued for resume.\n")
	done.ExpectStderr("")

	waitForStatus(t, flow, env, runID, file, "Succeeded")
	waitForFileContent(t, flow.ProjectPath("published.txt"), "published\n")
	waitForFileContent(t, flow.ProjectPath("attempts.txt"), ":\n1:add-tests\n")
}

func collapseStatusLines(stdout string) string {
	var result strings.Builder
	for line := range strings.SplitSeq(stdout, "\n") {
		if index := strings.Index(line, "│     "); index >= 0 {
			line = line[index+len("│     "):]
		}
		result.WriteString(strings.TrimSpace(line))
	}
	return result.String()
}

func runConcurrentCompletions(
	t *testing.T,
	flow *harness.Runner,
	env []string,
	runID string,
	first string,
	second string,
) []*harness.Result {
	t.Helper()
	values := []string{first, second}
	results := make([]*harness.Result, len(values))
	start := make(chan struct{})
	var wg sync.WaitGroup
	for index, value := range values {
		wg.Go(func() {
			<-start
			results[index] = complete(t, flow, env, runID, "choose", "concurrent.yaml", "--input=lane="+value)
		})
	}
	close(start)
	wg.Wait()
	return results
}

func TestArtifactReferencesResolveWhenTaskOpens(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-artifacts"

	startWaiting(t, flow, env, runID, "artifact_snapshot.yaml")
	status := waitForStatus(t, flow, env, runID, "artifact_snapshot.yaml", "Waiting")
	require.Contains(t, status.Stdout(), "reports/production/summary.txt")
	require.Contains(t, status.Stdout(), "reports/release.diff")
	require.NotContains(t, status.Stdout(), "${params.target}")

	// The run writes none of the referenced artifacts, and completion must not
	// care.
	result := complete(t, flow, env, runID, "review", "artifact_snapshot.yaml")
	result.ExpectExitCode(0)
	waitForStatus(t, flow, env, runID, "artifact_snapshot.yaml", "Succeeded")
}

func TestArtifactReferenceEscapingArtifactDirFailsWithoutOpening(t *testing.T) {
	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "spec031-artifact-escape"

	start := flow.RunWithEnv(env, "start", "--run-id="+runID,
		"--params=target=../../secret", "artifact_escape.yaml")
	start.ExpectNonZeroExitCode()
	start.ExpectStderrContains("human task artifacts", "parent directory segments")

	status := waitForStatus(t, flow, env, runID, "artifact_escape.yaml", "Failed")
	require.NotContains(t, status.Stdout(), "Waiting")
}
