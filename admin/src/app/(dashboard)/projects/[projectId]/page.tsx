import { ProjectWorkspace } from '@/features/projects/ui/project-workspace';

export default async function Page({ params }: { params: Promise<{ projectId: string }> }) {
  const { projectId } = await params;
  return <ProjectWorkspace projectId={projectId} />;
}
