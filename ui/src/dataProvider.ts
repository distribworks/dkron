import { fetchUtils } from 'ra-core';
import jsonServerProvider from 'ra-data-json-server';

export const apiUrl = window.DKRON_API_URL || 'http://localhost:8080/v1'
const dataProvider = jsonServerProvider(apiUrl);

const myDataProvider = {
    ...dataProvider,
    getManyReference: (resource: any, params: any) => {
        // Commented because executions query do not accept any query options yet
        
        // const { page, perPage } = params.pagination;
        // const { field, order } = params.sort;

        // const query = {
        //     ...fetchUtils.flattenObject(params.filter),
        //     [params.target]: params.id,
        //     _sort: field,
        //     _order: order,
        //     _start: (page - 1) * perPage,
        //     _end: page * perPage,
        // };
        const url = `${apiUrl}/${params.target}/${params.id}/${resource}`;

        return fetchUtils.fetchJson(url).then(({ headers, json }) => {
            if (!headers.has('x-total-count')) {
                throw new Error(
                    'The X-Total-Count header is missing in the HTTP Response. The jsonServer Data Provider expects responses for lists of resources to contain this header with the total number of results to build the pagination. If you are using CORS, did you declare X-Total-Count in the Access-Control-Expose-Headers header?'
                );
            }
            return {
                data: json,
                total: parseInt(
                    headers.get('x-total-count')!.split('/').pop() || '',
                    10
                ),
            };
        });
    }
}

export default myDataProvider;
