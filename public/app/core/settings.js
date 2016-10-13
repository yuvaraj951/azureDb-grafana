define([
  'lodash',
],
function (_) {
  "use strict";

  return function Settings (options) {
    var defaults = {
      datasources                   : {},
      window_title_prefix           : 'Icon - ',
      panels                        : {
      'clock-panel-master': { path: 'panels/clock-panel-master' },
      'singlestat': { path: 'panels/singlestat' }
      },
      new_panel_title: 'Panel Title',
      playlist_timespan: "1m",
      unsaved_changes_warning: true,
      appSubUrl: ""
    };

    return _.extend({}, defaults, options);
  };
});
