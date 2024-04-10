import React, { useEffect, useRef } from "react";
import * as echarts from "echarts";

interface CpuInfo {
  TimeStamps: string[];
  CpuUsages: number[];
}

const Cpu: React.FC<CpuInfo> = ({ TimeStamps, CpuUsages }) => {
  const chartRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (chartRef.current) {
      const myChart = echarts.init(chartRef.current);

      const option = {
        tooltip: {
          trigger: "axis",
          position: function (pt: any[]) {
            return [pt[0], "10%"];
          },
        },

        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: "none",
            },
            restore: {},
            saveAsImage: {},
          },
        },
        xAxis: {
          type: "category",
          boundaryGap: false,
          data: TimeStamps,
        },
        yAxis: {
          type: "value",
          boundaryGap: [0, "100%"],
        },
        dataZoom: [
          {
            type: "inside",
            start: 0,
            end: 10,
          },
          {
            start: 0,
            end: 10,
          },
        ],
        series: [
          {
            name: "Cpu Usage",
            type: "line",
            symbol: "none",
            sampling: "lttb",
            itemStyle: {
              color: "rgb(255, 70, 131)",
            },
            areaStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                  offset: 0,
                  color: "rgb(255, 158, 68)",
                },
                {
                  offset: 1,
                  color: "rgb(255, 70, 131)",
                },
              ]),
            },
            data: CpuUsages,
          },
        ],
      };

      myChart.setOption(option);

      // 清理函数
      return () => {
        myChart.dispose();
      };
    }
  }, [TimeStamps, CpuUsages]); // 当 TimeStamps 或 CpuUsages 更新时，重新渲染图表

  return <div ref={chartRef} style={{ width: "100%", height: "400px" }} />;
};

export default Cpu;
