import Ember from "ember";

export default Ember.Route.extend({
	titleToken: "Welcome",
	model: function() {
		return Ember.RSVP.hash({
			updates: Ember.$.getJSON("/api/v1/updates"),
			advisories: [],
			comments: []
		});
	}
});
