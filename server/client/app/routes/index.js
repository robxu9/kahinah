import Ember from "ember";

export default Ember.Route.extend({
  titleToken: "Welcome",
  model: function() {
    return Ember.RSVP.hash({
      updates: [],
      advisories: [],
      comments: []
    });
  }
});
