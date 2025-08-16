import React from "react";
import './App.css';
import { Route, Switch } from 'react-rooter-dom';
import Header from './components/Header/Header';
import ChatPage from "./components/Chat/Chatpage";
import LoginPage from './components/Login/LoginPage';
import { ProtectedRoute } from './authorization/protected.route';
