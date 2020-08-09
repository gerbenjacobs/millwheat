let queuePanes;

$(document).ready(function () {
    queuePanes = $(".queue_pane");
    barracksPanes = $(".barracks_pane");

    // Queue building
    $("#queue_tabs a").on("click", clickTabs(queuePanes));
    $("#barracks_tabs a").on("click", clickTabs(barracksPanes));

    // Town Building tabs
    $(".tb_tab").on("click", clickTBTabs);
});

function clickTabs(panes) {
    return function () {
        $(this).siblings().removeClass("active")
        $(this).addClass("active");

        panes.prop("hidden", true)
        let selected = $(this).attr("aria-controls");
        $("#" + selected).prop("hidden", false)
    }
}

function clickTBTabs() {
    $(this).siblings().removeClass("active")
    $(this).addClass("active");

    let selected = $(this).attr("aria-controls");
    $(".tb_pane_" + selected.split("_")[1]).prop("hidden", true)
    $("#" + selected).prop("hidden", false)
}