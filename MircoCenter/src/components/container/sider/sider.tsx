import { Layout } from "antd";
import Tree from "../../tree/tree.tsx";

const Sider = () => {
  return (
    <Layout.Sider theme={"light"}>
      <Tree />
    </Layout.Sider>
  );
};

export default Sider;
