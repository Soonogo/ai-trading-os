import { PageHeader } from '@/components/page-header';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@aitos/design-system';

export default function PromptStudioPage() {
  return (
    <div className="space-y-6">
      <PageHeader title="Prompt Studio" description="Version, test, and deploy agent prompts." />

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-medium">Overview</CardTitle>
            <CardDescription>Live data will appear here</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-24 rounded-lg bg-white/[0.03]" />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-medium">Activity</CardTitle>
            <CardDescription>Recent agent events</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-24 rounded-lg bg-white/[0.03]" />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base font-medium">Insights</CardTitle>
            <CardDescription>AI-generated takeaways</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-24 rounded-lg bg-white/[0.03]" />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
