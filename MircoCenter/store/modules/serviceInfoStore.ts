import { createSlice } from "@reduxjs/toolkit";
import { FetchServiceInfo } from "../../utils/api";

const serviceInfoStore = createSlice({
  name: "serviceInfo",
  initialState: {
    service_instances: [],
  },
  reducers: {
    setServiceInfo: (state, action) => {
      state.service_instances = action.payload;
    },
  },
});

const { setServiceInfo } = serviceInfoStore.actions;
export async function GetServiceInfo(dispatch: any) {
  const result = await FetchServiceInfo();
  dispatch(setServiceInfo(result.data.service_instances));
}

export const serviceInfoReducer = serviceInfoStore.reducer;
