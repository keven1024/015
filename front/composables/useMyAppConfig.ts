const useMyAppConfig = () => {
  const { data } = useFetch("/config");
  return data;
};

export default useMyAppConfig;
