import { ScenarioEditor } from '@/features/scenarios/ui/scenario-editor';

export default async function Page({
  params,
}: {
  params: Promise<{ projectId: string; scenarioId: string }>;
}) {
  const { projectId, scenarioId } = await params;
  return <ScenarioEditor projectId={projectId} scenarioId={scenarioId} />;
}
