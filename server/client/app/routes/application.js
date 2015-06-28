import Ember from "ember";

export default Ember.Route.extend({
  title: function(tokens) {
    return tokens.join(' - ') + ' • Kahinah';
  },
  actions: {
    error: function(error, transition) {
      this.set('previousTransition', transition);

      if (error.status === 0) {
        return this.transitionTo('errors.notconnected');
      } else if (error.status === 403) {
        return this.transitionTo('errors.notauthorized');
      } else if (error.status === 404) {
        return this.transitionTo('errors.notfound', {
          path: "not-found"
        });
      } else if (error.status === 500) {
        return this.transitionTo('errors.internal');
      }

      console.log(error);
      this.render('components.flash', {
        into: 'application',
        outlet: 'flash',
        model: {
          text: 'Unknown error occurred. Check the web console, or refresh to try again.',
          error: error,
          type: 'alert',
        }
      });

      return true;
    },
    refresh: function() {
      this.refresh();
    },
    retryPrevious: function() {
      var previousTransition = this.get('previousTransition');
      if (previousTransition) {
        previousTransition.retry();
      }
    }
  }
});
