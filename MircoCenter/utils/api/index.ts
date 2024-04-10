import { request } from "../../axios/axios";

export default interface Service {
  project_name: string;
  service_name: string;
  service_host: string;
}
// 获取服务状态
export async function FetchServiceStatus(service: Service) {
  const result = await request.get(
    `/ServiceManager/ServiceStatus?projectname=${service.project_name}&servicename=${service.service_name}&servicehost=${service.service_host}`,
  );
  return result.data;
}

// 获取服务信息
export async function FetchServiceInfo() {
  const result = await request.get("/ServiceManager");
  return result.data;
}

//  调整某服务优先级
export async function SetBestService(service: Service) {
  const result = await request.post("/ServiceManager/SetBestService", service);
  return result.data;
}

//  使能某服务不是最优
export async function UnSetBestService(service: Service) {
  const result = await request.delete("/ServiceManager/UnSetBestService", {
    data: service,
  });
  return result.data;
}

//  服务注销
export async function UnRegisterService(service: Service) {
  const result = await request.delete("/ServiceManager/UnRegisterService", {
    data: service,
  });
  return result.data;
}

// 获取指定服务两小时内的指标
export async function FetchServiceMetricsForTwoHour(service: Service) {
  const result = await request.get(
    `/ServiceManager/Metrics/TwoHour?projectname=${service.project_name}&servicename=${service.service_name}&servicehost=${service.service_host}`,
  );
  return result.data;
}

// 获取指定服务六小时内的指标
export async function FetchServiceMetricsForSixHour(service: Service) {
  const result = await request.get(
    `/ServiceManager/Metrics/SixHour?projectname=${service.project_name}&servicename=${service.service_name}&servicehost=${service.service_host}`,
  );
  return result.data;
}

// 获取指定服务12小时内的指标
export async function FetchServiceMetricsForTwelveHour(service: Service) {
  const result = await request.get(
    `/ServiceManager/Metrics/TwelveHour?projectname=${service.project_name}&servicename=${service.service_name}&servicehost=${service.service_host}`,
  );
  return result.data;
}

// 获取指定服务24小时内的指标
export async function FetchServiceMetricsForLastDay(service: Service) {
  const result = await request.get(
    `/ServiceManager/Metrics/LastDay?projectname=${service.project_name}&servicename=${service.service_name}&servicehost=${service.service_host}`,
  );
  return result.data;
}
