import { fetchUtils } from 'ra-core';
import jsonServerProvider from 'ra-data-json-server';
import { stringify } from 'query-string';

export const apiUrl = window.DKRON_API_URL || 'http://localhost:8080/v1'
const dataProvider = jsonServerProvider(apiUrl);

const myDataProvider = {
    ...dataProvider,
    getManyReference: (resource: any, params: any) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;

        const query = {
            ...fetchUtils.flattenObject(params.filter),
            [params.target]: params.id,
            _sort: field,
            _order: order,
            _start: (page - 1) * perPage,
            _end: page * perPage,
            output_size_limit: 200,
        };
        const url = `${apiUrl}/${params.target}/${params.id}/${resource}?${stringify(query)}`;

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
