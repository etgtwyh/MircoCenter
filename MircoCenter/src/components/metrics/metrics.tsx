import { useEffect, useState } from "react";
import { Button, Card, Col, Row, Switch, Tabs } from "antd";
import Cpu from "../cpu/cpu";
import Memory from "../memory/memory";
import Icmp from "../icmp/icmp";
import { useNavigate, useSearchParams } from "react-router-dom";
import {
  FetchServiceMetricsForTwoHour,
  FetchServiceMetricsForSixHour,
  FetchServiceMetricsForTwelveHour,
  FetchServiceMetricsForLastDay,
  SetBestService,
  UnSetBestService,
  FetchServiceStatus,
  UnRegisterService,
} from "../../../utils/api";

interface MetricsData {
  timestamps: string[];
  cpuUsages: number[];
  memoryUsages: number[];
  icmpDelays: number[];
}

const Metrics = () => {
  const [searchParams] = useSearchParams();
  const Navigate = useNavigate();
  const projectname = searchParams.get("projectname");
  const servicename = searchParams.get("servicename");
  const servicehost = searchParams.get("servicehost");
  const [metrics, setMetrics] = useState<MetricsData>({
    cpuUsages: [],
    icmpDelays: [],
    memoryUsages: [],
    timestamps: [],
  });
  const writeoff = () => {
    const wo = async () => {
      if (!projectname || !servicename || !servicehost) return;
      UnRegisterService({
        project_name: projectname,
        service_name: servicename,
        service_host: servicehost,
      });
      Navigate("/");
      window.location.reload();
    };
    wo();
  };
  // 使用状态控制Switch的选中状态
  const [isSwitchChecked, setIsSwitchChecked] = useState(false);

  const confirmServiceStatus = async () => {
    if (!projectname || !servicename || !servicehost) return;
    let result;
    result = await FetchServiceStatus({
      project_name: projectname,
      service_name: servicename,
      service_host: servicehost,
    });
    setIsSwitchChecked(result.data.priority);
  };
  const fetchMetrics = async (timeRange: string) => {
    if (!projectname || !servicename || !servicehost) return;

    let result;
    switch (timeRange) {
      case "2h":
        result = await FetchServiceMetricsForTwoHour({
          project_name: projectname,
          service_name: servicename,
          service_host: servicehost,
        });
        break;
      case "6h":
        result = await FetchServiceMetricsForSixHour({
          project_name: projectname,
          service_name: servicename,
          service_host: servicehost,
        });
        break;
      case "12h":
        result = await FetchServiceMetricsForTwelveHour({
          project_name: projectname,
          service_name: servicename,
          service_host: servicehost,
        });
        break;
      case "24h":
        result = await FetchServiceMetricsForLastDay({
          project_name: projectname,
          service_name: servicename,
          service_host: servicehost,
        });
        break;
      default:
        result = null;
    }

    if (result) {
      setMetrics(result.data);
    }
  };

  useEffect(() => {
    confirmServiceStatus();
    fetchMetrics("2h");
    // eslint-disable-next-line
  }, [projectname, servicename, servicehost, searchParams.toString()]); // Trigger re-fetch when search params change

  const onChange = (key: string) => {
    fetchMetrics(key);
  };

  const items = [
    { label: "2h", key: "2h" },
    { label: "6h", key: "6h" },
    { label: "12h", key: "12h" },
    { label: "24h", key: "24h" },
  ].map((item) => ({
    key: item.key,
    label: item.label,
    children: (
      <Row gutter={16}>
        <Col span={8}>
          <Card title="CPU Usage" bordered={false}>
            <Cpu
              CpuUsages={metrics.cpuUsages}
              TimeStamps={metrics.timestamps}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card title="Memory Usage" bordered={false}>
            <Memory
              MemoryUseages={metrics.memoryUsages}
              TimeStamps={metrics.timestamps}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card title="ICMP Delay" bordered={false}>
            <Icmp
              icmpDelays={metrics.icmpDelays}
              TimeStamps={metrics.timestamps}
            />
          </Card>
        </Col>
      </Row>
    ),
  }));

  return (
    <div className="mt-10 pl-10 pr-10">
      {projectname && servicename && servicehost && (
        <Switch
          checkedChildren="Set as Best"
          unCheckedChildren="Unset Best"
          checked={isSwitchChecked}
          defaultChecked={false}
          onChange={(checked, _event) => {
            if (checked) {
              SetBestService({
                project_name: projectname,
                service_name: servicename,
                service_host: servicehost,
              });
            } else {
              UnSetBestService({
                project_name: projectname,
                service_name: servicename,
                service_host: servicehost,
              });
            }
            setIsSwitchChecked(checked);
          }}
        />
      )}
      <Button danger shape={"round"} className={"ml-4"} onClick={writeoff}>
        write off
      </Button>
      <Tabs defaultActiveKey="2h" onChange={onChange} items={items} />
    </div>
  );
};

export default Metrics;
