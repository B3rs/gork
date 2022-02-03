const BASE_URL = window.location.href;

export const fetchJob = async (id) => {
  try {
    const resp = await fetch(`${BASE_URL}api/v1/jobs/${id}`);
    if (resp.status !== 200) {
      throw new Error(resp.statusText);
    }

    return resp.json();
  } catch (error) {
    console.log(error);
  }
};

export const createJob = async (params) => {
  try {
    const resp = await fetch(`${BASE_URL}api/v1/jobs/create`, {
      method: "POST",
    });
    if (resp.status !== 200) {
      throw new Error(resp.statusText);
    }
    return await resp.json();
  } catch (error) {
    console.log(error);
  }
};
