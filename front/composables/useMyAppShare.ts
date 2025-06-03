const downloadFile = async (share_id: string) => {
  const data = await $fetch<{
    code: number;
    data: {
      token?: string;
    };
  }>(`/api/download`, {
    method: "POST",
    body: {
      share_id,
    },
  });
  const { token } = data?.data || {};
  if (!token) {
    return;
  }
  (window as any)?.open(`/api/download?token=${token}`);
};

const createShare = async (data: any) => {
  return await $fetch<{
    code: number;
    data: {
      id?: string;
    };
  }>(`/api/share`, {
    method: "POST",
    body: data,
  });
};

const createFileShare = async (data: {
  file_id: string;
  config: {
    download_nums: number;
    expire_time: number;
    has_pickup_code?: boolean;
    has_password?: boolean;
    pickup_code?: string;
    password?: string;
    notify_email?: string;
  };
  file_name: string;
}) => {
  const { file_id, config, file_name } = data || {};
  return await createShare({
    type: "file",
    data: file_id,
    config,
    file_name,
  });
};

const createTextShare = async (data: { text: string; config: any }) => {
  const { text, config } = data || {};
  return await createShare({
    type: "text",
    data: text,
    config,
  });
};

const useMyAppShare = () => {
  return {
    downloadFile,
    createShare,
    createFileShare,
    createTextShare,
  };
};

export default useMyAppShare;
