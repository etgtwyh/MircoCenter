import { createHashRouter } from "react-router-dom";
import Layout from "../src/components/container/layout.tsx";
import Content from "../src/components/container/content/content.tsx";
import Metrics from "../src/components/metrics/metrics.tsx";
import EChartsTree from "../src/components/etree/EChartsTree.tsx";

export const router = createHashRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      {
        element: <Content />,
        children: [
          {
            index: true,
            element: <EChartsTree />,
          },
          {
            path: "/metrics",
            element: <Metrics />,
          },
        ],
      },
    ],
  },
]);
