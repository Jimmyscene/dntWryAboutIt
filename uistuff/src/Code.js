import React, { useState, useEffect } from 'react';
import Highlight from 'react-highlight'


class Code extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            "data": "Loading"
        }
        this.setData = this.setData.bind(this)
        this.getData = this.getData.bind(this)
    }
    setData(data) {
        return this.setState({data})
    }
    componentDidUpdate(prevProps) {
        if (prevProps.debug != this.props.debug) {
            return this.getData()
        }
    }

    componentDidMount() {
        return this.getData()
    }
    getData() {
        console.log(this.props)
        return fetch(`http://localhost:9000/file?debug=${this.props.debug}`, {mode: 'cors'})
            .then((resp) => resp.json())
            .then(({data}) => this.setData(data))
            .catch((err) => this.setData("Errored" + err))
    }
    render() {
        return (
            <Highlight className="python">
            {this.state.data}
            </Highlight>
        )
    }
}

export default Code
