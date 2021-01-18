import React from "react";
// import TrafficSection from "./TrafficSection";
import WikiSection from "./WikiSection";

class DashBoard extends React.Component {
    render() {
        return (
            <div className="dashboard">
                <div className="dashboard-board">
                    {/* <TrafficSection/> */}
                    <WikiSection />

                </div>
            </div>
        );
    }
}

// module.exports = DashBoard;
export default DashBoard;