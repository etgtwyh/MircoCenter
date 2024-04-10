import { Layout } from "antd";

import { Outlet } from "react-router-dom";

const Content = () => {
  return (
    <Layout.Content>
      <div className={"min-h-lvh min-w-full "}>
        <Outlet />
      </div>
    </Layout.Content>
  );
};

export default Content;
