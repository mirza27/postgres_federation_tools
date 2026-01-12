import { Outlet } from "react-router-dom";
import { Sidebar } from "@/components/sidebar";

export default function MainLayout() {
  return (
    <div className="flex h-screen bg-background">
      <Sidebar />
      <Outlet />
    </div>
  );
}
