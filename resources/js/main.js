let queueTabs, queuePanes;

$(document).ready(function () {
    queueTabs = $("#queue_tabs");
    queuePanes = $(".queue_pane");

    $("#queue_tabs a").on("click", clickTabs(queueTabs, queuePanes));
});

function clickTabs(tabs, panes) {
    return function() {
        console.log(tabs, panes)
        tabs.children("a").removeClass("active");
        $(this).addClass("active");

        panes.prop("hidden", true)
        let selected = $(this).attr("aria-controls");
        $("#" + selected).prop("hidden", false)
    }
}