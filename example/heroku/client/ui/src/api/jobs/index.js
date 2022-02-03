const BASE_URL = "https://gork-client-example.herokuapp.com/"; ///window.location.href;

export const fetchJob = async (id) => {
  try {
    const resp = await fetch(`${BASE_URL}api/v1/jobs/${id}`);
    return await resp.json();
  } catch (error) {
    console.log(error);
  }
};

export const createJob = async (params) => {
  try {
    const resp = await fetch(`${BASE_URL}api/v1/jobs`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    });
    return await resp.json();
  } catch (error) {
    console.log(error);
  }
};
