import { configureStore } from "@reduxjs/toolkit";
import { serviceInfoReducer } from "./modules/serviceInfoStore.ts";

export const store = configureStore({
  reducer: serviceInfoReducer,
});
