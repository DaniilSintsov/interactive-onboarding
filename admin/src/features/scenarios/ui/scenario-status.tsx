import type { ScenarioStatus as Status } from '@/shared/api/types';

const labels: Record<Status, string> = {
  in_development: 'В разработке',
  enabled: 'Активен',
  disabled: 'Выключен',
};

export function ScenarioStatus({ status }: { status: Status }) {
  return <span className={`status-pill status-${status}`}>{labels[status]}</span>;
}
