import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router";

import { InventoryPage } from "./routes/InventoryPage";
import "./styles.css";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  { path: "/", element: <InventoryPage /> },
  { path: "*", element: <InventoryPage /> },
]);

const root = document.getElementById("root");
if (!root) {
  throw new Error("root element is missing");
}

createRoot(root).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
);
