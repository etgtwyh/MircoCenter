import Icon from "@ant-design/icons";

import type { GetProps } from "antd";

type CustomIconComponentProps = GetProps<typeof Icon>;

const OnlineSvg = () => (
  <svg
    d="1712583141135"
    className="icon"
    viewBox="0 0 1024 1024"
    version="1.1"
    xmlns="http://www.w3.org/2000/svg"
    p-id="5608"
    width="20"
    height="20"
  >
    <path
      d="M539.946667 1000.064C805.76 769.664 938.666667 578.517333 938.666667 426.666667c0-235.648-191.018667-426.666667-426.666667-426.666667S85.333333 191.018667 85.333333 426.666667c0 151.893333 132.906667 342.997333 398.72 573.397333a42.666667 42.666667 0 0 0 55.893334 0z"
      fill="#3CC05C"
      p-id="5609"
    ></path>
    <path
      d="M512 217.173333a209.493333 209.493333 0 1 0 0 418.986667 209.493333 209.493333 0 0 0 0-418.986667zM512 170.666667a256 256 0 1 1 0 512 256 256 0 0 1 0-512z m99.328 174.592a23.253333 23.253333 0 1 1 32.896 32.896l-150.613333 150.656-114.432-114.389334a23.253333 23.253333 0 0 1 32.896-32.938666l81.493333 81.493333 117.76-117.76z"
      fill="#FFFFFF"
      p-id="5610"
    ></path>
  </svg>
);

export const OnlineIcon = (props: Partial<CustomIconComponentProps>) => (
  <Icon component={OnlineSvg} {...props} />
);
