import React from "react";
import TrafficSection from "./TrafficSection";

class DashBoard extends React.Component {
    render() {
        return (
            <div className="dashboard">
                <div className="dashboard-board">
                    <TrafficSection/>
                </div>
            </div>
        );
    }
}

// module.exports = DashBoard;
export default DashBoard;