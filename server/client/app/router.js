import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
	location: config.locationType
});

Router.map(function() {
	this.route('updates', {
		path: '/updates'
	});
	this.resource('update', {
		path: '/updates/:id'
	});
	// errors
	this.route('errors.internal', {
		path: '/monkeys'
	});
	this.route('errors.notauthorized', {
		path: '/not-authorized'
	});
	this.route('errors.notconnected', {
		path: '/no-interwebs'
	});
	this.route('errors.notfound', {
		path: '/*path'
	});
});

export default Router;
