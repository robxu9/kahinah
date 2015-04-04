import Ember from "ember";

export default Ember.Route.extend({
	titleToken: "Latest Updates",
	model: function() {
		return Ember.$.getJSON("/api/v1/updates");
	}
});
