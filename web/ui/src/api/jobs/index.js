export const fetchJobs = async (opts = {limit: 100, page: 0, search: ""})=> {
    const {limit, page, search} = opts;
    try{
        const resp = await fetch(`http://localhost:8080/api/v1/jobs?limit=${limit}&page=${page}&q=${search}`);
        if (resp.status !== 200) {
            throw new Error(resp.statusText);
        }
    
        return resp.json();
    }  catch (error) {
        console.log(error);
    }
}

export const fetchJob = async (id)=> {
    try{
        const resp = await fetch(`http://localhost:8080/api/v1/jobs/${id}`);
        if (resp.status !== 200) {
            throw new Error(resp.statusText);
        }
    
        return resp.json();
    }  catch (error) {
        console.log(error);
    }
}
  
export const retryJob = async (id) => {
    try{
        const resp = await fetch(`http://localhost:8080/api/v1/jobs/${id}/retry`,{method: 'POST'});
        if (resp.status !== 200) {
            throw new Error(resp.statusText);
        }
    }  catch (error) {
        console.log(error);
    }
}