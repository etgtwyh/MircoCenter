import React, { useEffect, useRef } from "react";
import { useSelector } from "react-redux";
import * as echarts from "echarts";
import type { EChartsOption } from "echarts";

type ServiceInstances = {
  [projectName: string]: {
    [serviceName: string]: string[]; // 这里假设服务相关的字符串直接放在数组中
  };
}[];

interface AllServices {
  service_instances: ServiceInstances;
}
// ECharts 树节点接口
interface EChartsTreeNode {
  name: string;
  children?: EChartsTreeNode[];
}
const convertServiceInstancesToEChartsTreeData = (
  serviceInstances: ServiceInstances,
): EChartsTreeNode[] => {
  // 创建一个根节点
  const rootNode: EChartsTreeNode = {
    name: "ServiceCenter",
    children: [],
  };

  // 将每个服务实例添加为根节点的子节点
  if (serviceInstances != null) {
    serviceInstances.forEach((instance) => {
      Object.entries(instance).forEach(([projectName, services]) => {
        // 创建每个项目的节点
        const projectNode: EChartsTreeNode = {
          name: projectName,
          children: [],
        };

        Object.entries(services).forEach(([serviceName, hosts]) => {
          // 确保 hosts 是数组
          if (Array.isArray(hosts)) {
            // 创建每个服务的节点，并将主机作为子节点添加
            const serviceNode: EChartsTreeNode = {
              name: serviceName,
              children: hosts.map((host) => ({ name: host })),
            };
            projectNode.children?.push(serviceNode);
          } else {
            console.error(
              `Hosts for service "${serviceName}" is not an array:`,
              hosts,
            );
          }
        });

        rootNode.children?.push(projectNode);
      });
    });
  }
  // 返回包含根节点的数组
  return [rootNode];
};

// 假设我们的 Redux store 有一个名为 allServices 的部分，并且它符合 AllServices 接口
const EChartsTree: React.FC = () => {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const serviceInstances = useSelector(
    (state: AllServices) => state.service_instances,
  );

  useEffect(() => {
    if (chartRef.current) {
      const myChart = echarts.init(chartRef.current);
      const treeData =
        convertServiceInstancesToEChartsTreeData(serviceInstances);

      const option: EChartsOption = {
        tooltip: {
          trigger: "item",
          triggerOn: "mousemove",
        },
        series: [
          {
            type: "tree",
            data: treeData,
            top: "50%",
            left: "10%",
            bottom: "1%",
            right: "10%",
            symbolSize: 7,
            label: {
              position: "left",
              verticalAlign: "middle",
              align: "right",
              fontSize: 12,
            },
            leaves: {
              label: {
                position: "right",
                verticalAlign: "middle",
                align: "left",
                fontSize: 12,
              },
            },
            emphasis: {
              focus: "descendant",
            },
            expandAndCollapse: true,
            animationDuration: 550,
            animationDurationUpdate: 750,
            lineStyle: {
              color: "#333333", // 设置线条的颜色，这里为黑色
              width: 1, // 设置线条的宽度
              // 也可以添加其他样式属性，如线型：'solid', 'dashed', 'dotted'
            },
          },
        ],
      };

      myChart.setOption(option);

      return () => {
        myChart.dispose(); // 组件卸载时清理 ECharts 实例
      };
    }
  }, [serviceInstances]); // 当 serviceInstances 更新时重新渲染图表

  return <div ref={chartRef} style={{ width: "100%", height: "500px" }} />;
};

export default EChartsTree;
