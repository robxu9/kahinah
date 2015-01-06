import Ember from "ember";
import DS from "ember-data";

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
