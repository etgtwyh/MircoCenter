import { Layout as AntdLayout } from "antd";
import Sider from "./sider/sider";

import { Outlet } from "react-router-dom";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { GetServiceInfo } from "../../../store/modules/serviceInfoStore.ts";

const Layout = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    GetServiceInfo(dispatch);
  }, [dispatch]);
  return (
    <AntdLayout style={{ minHeight: "100vh" }}>
      <Sider />
      <Outlet />
    </AntdLayout>
  );
};

export default Layout;
