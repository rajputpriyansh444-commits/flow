import type { ApiDAGDetails, ApiDAGRunDetails } from '../../contracts';
import type { WorkspaceRunValue } from '../../data';

export interface DockTabBodyProps {
  run: WorkspaceRunValue;
  fileName: string;
  dag: ApiDAGDetails;
  /**
   * The host shows this step: the Log tab opens its log, whatever the URL
   * says. Without it the URL and the auto-open choose.
   */
  step?: string;
}

/** Props of a run tab body after the frame has checked that details exist. */
export interface RunTabProps extends DockTabBodyProps {
  details: ApiDAGRunDetails;
}
