import React, { useEffect, useRef } from "react";
import * as echarts from "echarts";
import type { EChartsOption } from "echarts";

interface IcmpDelay {
  TimeStamps: string[];
  icmpDelays: number[]; // ICMP延迟值，以毫秒为单位
}

const Icmp: React.FC<IcmpDelay> = ({ TimeStamps = [], icmpDelays = [] }) => {
  const chartRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (chartRef.current) {
      const myChart = echarts.init(chartRef.current);

      // 将 ICMP 延迟值从毫秒转换为秒
      const delaysInSeconds = icmpDelays.map((delay) => delay / 10000);

      const option: EChartsOption = {
        tooltip: {
          trigger: "axis",
          position: function (pt) {
            return [pt[0], "10%"];
          },
        },
        title: {
          left: "center",
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
          axisLabel: {
            formatter: "{value} s",
          },
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
            name: "ICMP Delay",
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
            data: delaysInSeconds,
          },
        ],
      };

      myChart.setOption(option);

      return () => {
        myChart.dispose();
      };
    }
  }, [TimeStamps, icmpDelays]); // 当 TimeStamps 或 icmpDelays 更新时，重新渲染图表

  return <div ref={chartRef} style={{ width: "100%", height: "400px" }} />;
};

export default Icmp;
