import React from 'react';
import { render, screen } from '@testing-library/react';

import '@testing-library/jest-dom'

import Form from '../Form'
import MessageBox from '../MessageBox';
import MessageStream from '../MessageStream';

describe('Individual component render testing', () => {

  beforeEach(() => {
    // write someting before each test
  });

  afterEach(() => {
    // write someting after each test
  });

  it('Renders a Form ', async () => {
    render(
      <Form />
    );
    const sendButton = screen.getByText('Send');

    expect(sendButton).toBeInTheDocument();
  });

  it('Renders a MessageBox', async () => {
    render(
      <MessageBox  message={{id: 0, username: "Username", message: "hello world !"}} />
    );
    const username = screen.getByText(/Username/)
    const message = screen.getByText(/hello world !/);

    expect(username).toBeInTheDocument();
    expect(message).toBeInTheDocument();
  });

  it('Renders a MessageStream', async () => {
    const messages = [{
      id: 0,
      username: "Alice",
      message: "Hello world !"
    }, {
      id: 1,
      username: "Bob",
      message: "I am a test !"
    }]
    render(
      <MessageStream messages={messages}/>
    )
    
    messages.forEach((item) => {
      const username = screen.getByText(new RegExp(`${item.username}`))
      const message = screen.getByText(new RegExp(`${item.message}`));

      expect(username).toBeInTheDocument();
      expect(message).toBeInTheDocument();
    })
  })
});