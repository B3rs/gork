const BASE_URL =
  process.env.NODE_ENV == "development"
    ? "http://localhost:8080/"
    : window.location.href;

export const fetchJobs = async (opts = { limit: 100, page: 0, search: "" }) => {
  const { limit, page, search } = opts;
  try {
    const resp = await fetch(
      `${BASE_URL}api/v1/jobs?limit=${limit}&page=${page}&q=${search}`
    );
    if (resp.status !== 200) {
      throw new Error(resp.statusText);
    }

    return resp.json();
  } catch (error) {
    console.log(error);
  }
};

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

export const retryJob = async (id) => {
  try {
    const resp = await fetch(`${BASE_URL}api/v1/jobs/${id}/retry`, {
      method: "POST",
    });
    if (resp.status !== 200) {
      throw new Error(resp.statusText);
    }
  } catch (error) {
    console.log(error);
  }
};

export const cancelJob = async (id) => {
  try {
    const resp = await fetch(`${BASE_URL}api/v1/jobs/${id}`, {
      method: "DELETE",
    });
    if (resp.status !== 200) {
      throw new Error(resp.statusText);
    }
  } catch (error) {
    console.log(error);
  }
};
