import React, { Component } from 'react';
import Logo  from './logo';
import './App.css';
import Checkbox from '@material-ui/core/Checkbox';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';


class App extends Component {


  render() {
    const theme = createMuiTheme({
      palette: {
        type: 'dark', // Switching the dark mode on is a single property value change.
      },
      typography: { useNextVariants: true },
    });


    return (
      <div className="App">
      <MuiThemeProvider theme={theme}>
      <div className="App-header">
            {Logo("#61DAFB")}
            <FormGroup row>
              <FormControlLabel
                control={
                  <Checkbox variant="contained" checked={this.props.debug} onChange={(e) => this.props.debugHandler(e.target.checked)}/>
                }
                label="Debug Mode?"
              />
            </FormGroup>
            <FormGroup row>
              <FormControlLabel
                control={
                  <Button color="primary" variant="contained" onClick={(e) => this.props.initailizedHandler(true)}>
                  Run!
                  </Button>
                }
              />
            </FormGroup>
          </div>
      </MuiThemeProvider>
      </div>
    );
  }
}

export default App;
