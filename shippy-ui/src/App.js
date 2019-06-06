import React, { Component } from 'react';
import './App.css';
import CreateConsignment from './CreateConsignment';

class App extends Component {

    state = {
        authenticated: false,
        email: '',
        password: '',
        err: '',
    }

    login = () => {
        fetch(`http://localhost:3000/rpc`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                request: {
                    email: this.state.email,
                    password: this.state.password,
                },
                service: 'Ethan.MicroServicePractice.user',
                method: 'UserService.Auth',
            }),
        })
            .then(res => res.json())
            .then(res => {
                this.setState({
                    token: res.token,
                    authenticated: true,
                });
            })
            .catch(err => this.setState({ err, authenticated: false, }));
    }

    signup = () => {
        fetch(`http://localhost:3000/rpc`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                request: {
                    email: this.state.email,
                    password: this.state.password,
                    name: this.state.name,
                },
                service: 'Ethan.MicroServicePractice.user',
                method: 'UserService.Create',
            }),
        })
            .then((res) => res.json())
            .then((res) => {
                this.setState({
                    token: res.token.token,
                    authenticated: true,
                });
                localStorage.setItem('token', res.token.token);
            })
            .catch(err => this.setState({ err, authenticated: false, }));
    }


    setEmail = e => {
        this.setState({
            email: e.target.value,
        });
    }

    setPassword = e => {
        this.setState({
            password: e.target.value,
        });
    }

    setName = e => {
        this.setState({
            name: e.target.value,
        });
    }

    renderLogin = () => {
        return (
            <div className="App-intro container">
                <br />
                <div className='Login'>
                    <div className='form-group'>
                        <input
                            type="email"
                            onChange={this.setEmail}
                            placeholder='E-Mail'
                            className='form-control' />
                    </div>
                    <div className='form-group'>
                        <input
                            type="password"
                            onChange={this.setPassword}
                            placeholder='Password'
                            className='form-control' />
                    </div>
                    <button className='btn btn-primary' onClick={this.login}>Login</button>
                    <br /><br />
                </div>
                <div className='Sign-up'>
                    <div className='form-group'>
                        <input
                            type='input'
                            onChange={this.setName}
                            placeholder='Name'
                            className='form-control' />
                    </div>
                    <div className='form-group'>
                        <input
                            type='email'
                            onChange={this.setEmail}
                            placeholder='E-Mail'
                            className='form-control' />
                    </div>
                    <div className='form-group'>
                        <input
                            type='password'
                            onChange={this.setPassword}
                            placeholder='Password'
                            className='form-control' />
                    </div>
                    <button className='btn btn-primary' onClick={this.signup}>Sign-up</button>
                </div>
            </div>
        );
    }

    renderAuthenticated = () => {
        return (
            <CreateConsignment token={this.state.token} />
        );
    }

    setToken = (token) => {
        localStorage.setItem('token', token);
    }

    getToken = () => {
        return localStorage.getItem('token');
    }

    isAuthenticated = () => {
        return this.state.token || this.getToken() || false;
    }

    render() {
        const authenticated = this.isAuthenticated();
        return (
            <div className="App">
                <div className="App-header">
                    <h2>Shippy</h2>
                </div>
                <div>
                    {(authenticated ? this.renderAuthenticated() : this.renderLogin())}
                </div>
            </div>
        );
    }
}

export default App;
