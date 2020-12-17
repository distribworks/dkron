export const darkTheme = {
    palette: {
        primary: {
            main: '#90caf9',
        },
        secondary: {
            main: '#FBBA72',
        },
        type: 'dark' as 'dark', // Switching the dark mode on is a single property value change.
    },
    overrides: {
        MuiAppBar: {
            colorSecondary: {
                color: '#ffffffb3',
                backgroundColor: '#616161e6',
            },
        },
    },
};

export const lightTheme = {
    palette: {
        primary: {
            main: '#4f3cc9',
        },
        secondary: {
            light: '#5f5fc4',
            main: '#283593',
            dark: '#001064',
            contrastText: '#fff',
        },
        background: {
            default: '#fcfcfe',
        },
        type: 'light' as 'light',
    },
    shape: {
        borderRadius: 10,
    },
    overrides: {
        RaMenuItemLink: {
            root: {
                borderLeft: '3px solid #fff',
            },
            active: {
                borderLeft: '3px solid #4f3cc9',
            },
        },
        MuiPaper: {
            elevation1: {
                boxShadow: 'none',
            },
            root: {
                border: '1px solid #e0e0e3',
                backgroundClip: 'padding-box',
            },
        },
        MuiButton: {
            contained: {
                backgroundColor: '#fff',
                color: '#4f3cc9',
                boxShadow: 'none',
            },
        },
        MuiAppBar: {
            colorSecondary: {
                color: '#ffffffb3',
                backgroundColor: '#333',
            },
        },
        MuiLinearProgress: {
            colorPrimary: {
                backgroundColor: '#f5f5f5',
            },
            barColorPrimary: {
                backgroundColor: '#d7d7d7',
            },
        },
        MuiFilledInput: {
            root: {
                backgroundColor: 'rgba(0, 0, 0, 0.04)',
                '&$disabled': {
                    backgroundColor: 'rgba(0, 0, 0, 0.04)',
                },
            },
        },
    },
};
