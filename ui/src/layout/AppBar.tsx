import { forwardRef } from 'react';
import { styled } from '@mui/material/styles';
import { AppBar, UserMenu, MenuItemLink, Link } from 'react-admin';
import Typography from '@mui/material/Typography';
import SettingsIcon from '@mui/icons-material/Settings';
import BookIcon from '@mui/icons-material/Book';
import Clock from './Clock';

import logo from '../images/dkron-logo.png';

const PREFIX = 'CustomAppBar';

const classes = {
    title: `${PREFIX}-title`,
    spacer: `${PREFIX}-spacer`,
    logo: `${PREFIX}-logo`
};

const StyledAppBar = styled(AppBar)({
    [`& .${classes.title}`]: {
        flex: 1,
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap',
        overflow: 'hidden',
    },
    [`& .${classes.spacer}`]: {
        flex: 1,
    },
    [`& .${classes.logo}`]: {
        maxWidth: "125px"
    },
});

const ConfigurationMenu = forwardRef<any, any>((props, ref) => {
    return (
        <MenuItemLink
            ref={ref}
            to="/settings"
            primaryText='Settings'
            leftIcon={<SettingsIcon />}
            onClick={props.onClick}
        />
    );
});

const CustomUserMenu = (props: any) => (
    <UserMenu {...props}>
        <MenuItemLink
            to="https://dkron.io/docs/basics/getting-started"
            primaryText='Docs'
            leftIcon={<BookIcon />}
        />
        {/* <ConfigurationMenu /> */}
    </UserMenu>
);

const CustomAppBar = (props: any) => {

    return (
        <StyledAppBar {...props} elevation={1} userMenu={<CustomUserMenu />}>
            <Typography
                variant="h6"
                color="inherit"
                className={classes.title}
                id="react-admin-title"
            />
            <div>
                <Link to="/" color="inherit" underline="none">
                    <img src={logo} alt="logo" className={classes.logo} />
                </Link>
            </div>
            <span className={classes.spacer} />
            <Clock/>
        </StyledAppBar>
    );
};

export default CustomAppBar;
