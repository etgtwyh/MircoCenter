import React from "react";
import { Spin, Tree as AntdTree } from "antd";
import { DownOutlined } from "@ant-design/icons";
import "./index.css";

import { useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";

type ServiceInstances = {
  [projectName: string]: {
    [serviceName: string]: string[]; // 这里假设服务相关的字符串直接放在数组中
  };
}[];

interface AllServices {
  service_instances: ServiceInstances;
}
const Tree: React.FC = () => {
  const Navigate = useNavigate();
  const service_instances = useSelector(
    (state: AllServices) => state.service_instances,
  );
  // onSelect事件处理器
  const onSelect = (selectedKeys: React.Key[]) => {
    if (typeof selectedKeys[0] === "string") {
      const keyParts = selectedKeys[0].split("/");
      if (keyParts.length === 3) {
        Navigate(
          `/metrics?projectname=${keyParts[0]}&servicename=${keyParts[1]}&servicehost=${keyParts[2]}`,
        );
      }
    }
  };

  // 动态生成treeData
  const treeData =
    service_instances !== null
      ? service_instances.flatMap((project, _projectIndex) =>
          Object.entries(project).map(([projectName, services]) => ({
            title: projectName,
            key: `root/${projectName}`,
            children: Object.entries(services).map(
              ([serviceName, serviceHosts], _serviceIndex) => ({
                title: serviceName,
                key: `${projectName}/${serviceName}`,
                children: serviceHosts.map((host, _hostIndex) => ({
                  title: host,
                  key: `${projectName}/${serviceName}/${host}`,
                  icon: <div className={"breathe"} />, // 使用自定义的在线状态图标
                })),
              }),
            ),
          })),
        )
      : null;

  if (treeData === null) {
    return <Spin />; // 渲染加载中组件
  } else {
    return (
      <AntdTree
        showIcon
        autoExpandParent={true}
        defaultSelectedKeys={["0-0-0"]}
        switcherIcon={<DownOutlined />}
        treeData={treeData}
        onSelect={onSelect} // 设置onSelect处理器
      />
    );
  }
};

export default Tree;
