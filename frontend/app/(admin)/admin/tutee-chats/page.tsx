
export default function AdminPage() {
  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold text-primary tracking-tight">Admin Dashboard</h2>
        <p className="text-foreground">
          Welcome to the admin dashboard.
        </p>
      </div>
      <div className="border-t">
        <div className="bg-background">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 p-6">
            <div className="rounded-xl border bg-card p-6">
              <h3 className="font-semibold">Welcome to Admin Panel</h3>
              <p className="text-muted-foreground mt-2">
                This is a protected admin area.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
} 