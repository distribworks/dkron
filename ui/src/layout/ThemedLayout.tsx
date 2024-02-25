import { Layout, LayoutProps, Sidebar } from 'react-admin';
import AppBar from './AppBar';

const CustomSidebar = (props: any) => <Sidebar {...props} size={200} />;

const ThemedLayout = (props: LayoutProps) => {
    return (
        <Layout
            {...props}
            appBar={AppBar}
            sidebar={CustomSidebar}
        />
    );
};
export default ThemedLayout;
