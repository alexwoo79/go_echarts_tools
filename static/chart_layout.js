(function () {
  function isPanelFullscreen(panel) {
    if (!panel) return false;
    return document.fullscreenElement === panel
      || document.webkitFullscreenElement === panel
      || document.mozFullScreenElement === panel
      || document.msFullscreenElement === panel;
  }

  function togglePanelFullscreen(panel) {
    if (!panel) return;
    if (isPanelFullscreen(panel)) {
      if (document.exitFullscreen) {
        document.exitFullscreen();
      } else if (document.webkitExitFullscreen) {
        document.webkitExitFullscreen();
      } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen();
      } else if (document.msExitFullscreen) {
        document.msExitFullscreen();
      }
      return;
    }

    if (panel.requestFullscreen) {
      panel.requestFullscreen();
    } else if (panel.webkitRequestFullscreen) {
      panel.webkitRequestFullscreen();
    } else if (panel.mozRequestFullScreen) {
      panel.mozRequestFullScreen();
    } else if (panel.msRequestFullscreen) {
      panel.msRequestFullscreen();
    }
  }

  function getChartTargetHeight(panel, opts) {
    if (!panel) return 0;
    var options = opts || {};
    var titleSelector = options.titleSelector || ".panel-title-row";
    var viewportBottomGap = options.viewportBottomGap != null ? options.viewportBottomGap : 14;
    var minHeight = options.minHeight != null ? options.minHeight : 420;

    var titleRow = panel.querySelector(titleSelector);
    var titleH = titleRow ? titleRow.offsetHeight : 52;
    var style = window.getComputedStyle(panel);
    var padTop = parseFloat(style.paddingTop) || 0;
    var padBottom = parseFloat(style.paddingBottom) || 0;

    var height;
    if (isPanelFullscreen(panel)) {
      height = window.innerHeight - titleH - padTop - padBottom;
    } else {
      var panelTop = panel.getBoundingClientRect().top;
      height = window.innerHeight - panelTop - viewportBottomGap - titleH - padTop - padBottom;
    }

    return Math.max(minHeight, Math.floor(height));
  }

  function syncChartHeight(panel, chartNode, opts) {
    if (!panel || !chartNode) return 0;
    var target = getChartTargetHeight(panel, opts);
    if (target > 0) {
      chartNode.style.height = target + "px";
    }
    return target;
  }

  window.ChartLayoutUtils = {
    isPanelFullscreen: isPanelFullscreen,
    togglePanelFullscreen: togglePanelFullscreen,
    getChartTargetHeight: getChartTargetHeight,
    syncChartHeight: syncChartHeight
  };
})();
