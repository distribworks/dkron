import { applyMiddleware, combineReducers, compose, createStore } from 'redux';
import { routerMiddleware, connectRouter } from 'connected-react-router';
import createSagaMiddleware from 'redux-saga';
import { all, fork } from 'redux-saga/effects';
import {
    adminReducer,
    adminSaga,
    USER_LOGOUT,
} from 'react-admin';
import storage from 'redux-persist/lib/storage';
import { persistReducer } from 'redux-persist';
import themeReducer from './themeReducer';

const createAdminStore = ({
    authProvider,
    dataProvider,
    history,
}) => {
    const persistConfig = {
        key: "dkronui",
        storage: storage,
        whitelist: ['theme']
    };
    const reducer = combineReducers({
        admin: adminReducer,
        router: connectRouter(history),
        theme: themeReducer,
        // our own reducers here
    });
    const persistedReducer = persistReducer(persistConfig, reducer);
    const resettableAppReducer = (state, action) =>
        persistedReducer(action.type !== USER_LOGOUT ? state : undefined, action);

    const saga = function* rootSaga() {
        yield all(
            [
                adminSaga(dataProvider, authProvider),
                // our own sagas here
            ].map(fork)
        );
    };
    const sagaMiddleware = createSagaMiddleware();

    const composeEnhancers =
        (process.env.NODE_ENV === 'development' &&
            typeof window !== 'undefined' &&
            window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ &&
            window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__({
                trace: true,
                traceLimit: 25,
            })) ||
        compose;
  
    const store = createStore(
        resettableAppReducer,
        { /* initial state here */ },
        composeEnhancers(
            applyMiddleware(
                sagaMiddleware,
                routerMiddleware(history),
                // our own middlewares here
            ),
            // our own enhancers here
        ),        
    );
    sagaMiddleware.run(saga);
    return store;
};

export default createAdminStore;
