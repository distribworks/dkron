import React, {Component} from 'react';

class Clock extends Component<{}, { date: Date }> {
    timer: any;

    constructor(props: any){
        super(props);
        this.state = {date: new Date()};
    }

    componentDidMount() {
        this.timer = setInterval(
            () => this.setState({date: new Date()}),
            1000
        );
    }

    componentWillUnmount() {
        clearInterval(this.timer);
    }

    render(){
        return( 
            <div className="clock">
                <div>{this.state.date.toLocaleTimeString()}</div>
            </div>
        )
    }
}
export default Clock
