import logging
from typing import List

# Ensure the Milvus_DB and ChunkedItem classes are accessible
from cognix_lib.db.milvus_db import Milvus_DB

from lib.spider.chunked_item import ChunkedItem

# Configure logging
logging.basicConfig(level=logging.DEBUG)

def test_store_chunk_list():
    # Initialize the Milvus_DB instance
    milvus_db = Milvus_DB()

    # Create a chunk list that will exceed the JSON field limit
    large_content = "A" * 70000  # This creates a string of 70,000 characters, exceeding the 65,536 limit

    large_content = ''' : automatic, 73
    2024 - 07 - 0
    9
    T11: 14:45.581915007
    Z
    Windows
    Forms
    layout
    model, 62
    2024 - 07 - 0
    9
    T11: 14:45.581915924
    Z
    WrapPanel, 64, 76
    2024 - 07 - 0
    9
    T11: 14:45.581916757
    Z
    WrapPanel

    class , 102

    2024 - 07 - 0
    9
    T11: 14:45.581917716
    Z
    layout, input, focus, and events(LIFE), 15
    2024 - 07 - 0
    9
    T11: 14:45.581918591
    Z
    layout
    panels, adding
    VisualStateManager
    to
    2024 - 07 - 0
    9
    T11: 14:45.581919507
    Z
    templates, 575
    2024 - 07 - 0
    9
    T11: 14:45.581920341
    Z
    layout
    pass, 583
    2024 - 07 - 0
    9
    T11: 14:45.581921216
    Z
    LayoutTransform
    property, 369, 468, 881, 993—
    2024 - 07 - 0
    9
    T11: 14:45.581922132
    Z
    994
    2024 - 07 - 0
    9
    T11: 14:45.581922966
    Z
    LCD
    monitors, native
    resolution, 8
    2024 - 07 - 0
    9
    T11: 14:45.581923841
    Z
    Left
    property, 94, 752, 756
    2024 - 07 - 0
    9
    T11: 14:45.581924757
    Z
    LIFE(layout, input, focus, and events), 15
    2024 - 07 - 0
    9
    T11: 14:45.581925632
    Z
    lifetime
    events, 133—134
    2024 - 07 - 0
    9
    T11: 14:45.581926507
    Z
    Light

    class , 896

    2024 - 07 - 0
    9
    T11: 14:45.581927341
    Z
    light
    sources, 896
    2024 - 07 - 0
    9
    T11: 14:45.581928216
    Z
    LightWave, 910
    2024 - 07 - 0
    9
    T11: 14:45.581929132
    Z
    line
    caps, in Lines and Polylines, 348
    2024 - 07 - 0
    9
    T11: 14:45.581930132
    Z
    Line

    class , 335

    2024 - 07 - 0
    9
    T11: 14:45.581930966
    Z
    inability
    to
    use
    flow
    content
    model, 344
    2024 - 07 - 0
    9
    T11: 14:45.581932007
    Z
    placing
    Line in Canvas, 344
    2024 - 07 - 0
    9
    T11: 14:45.581933091
    Z
    setting
    starting and ending
    points, 343
    2024 - 07 - 0
    9
    T11: 14:45.581934007
    Z
    Stroke
    property, 343
    2024 - 07 - 0
    9
    T11: 14:45.581934841
    Z
    understanding
    line
    caps, 348
    2024 - 07 - 0
    9
    T11: 14:45.581935757
    Z
    using
    negative
    coordinates
    for line, 343
        2024 - 07 - 0
        9
        T11: 14:45.581936632
        Z
    2024 - 07 - 0
    9
    T11: 14:45.581937424
    Z
    2024 - 07 - 0
    9
    T11: 14:45.581938257
    Z
    using
    StartLineCap and EndLineCap
    2024 - 07 - 0
    9
    T11: 14:45.581939132
    Z
    properties, 348
    2024 - 07 - 0
    9
    T11: 14:45.581940049
    Z
    line
    joins, 349
    2024 - 07 - 0
    9
    T11: 14:45.581940882
    Z
    linear
    interpolation
    2024 - 07 - 0
    9
    T11: 14:45.581941757
    Z
    animating
    property
    with special value of
    2024 - 07 - 0
    9
    T11: 14:45.581942632
    Z
    Double.NaN, 431
    2024 - 07 - 0
    9
    T11: 14:45.581945216
    Z
    animating
    two
    properties
    simultaneously,
    2024 - 07 - 0
    9
    T11: 14:45.581946091
    Z
    434
    2024 - 07 - 0
    9
    T11: 14:45.581947049
    Z
    Canvas as most
    common
    layout
    container
    2024 - 07 - 0
    9
    T11: 14:45.581947966
    Z
    for animation, 431
        2024 - 07 - 0
        9
        T11: 14:45.581948841
        Z
        creating
        additive
        animation
        by
        setting
    2024 - 07 - 0
    9
    T11: 14:45.581949716
    Z
    IsAdditive
    property, 433
    2024 - 07 - 0
    9
    T11: 14:45.581950591
    Z
    creating
    animation
    that
    widens
    button, 430
    2024 - 07 - 0
    9
    T11: 14:45.581951507
    Z
    description, 453
    2024 - 07 - 0
    9
    T11: 14:45.581952341
    Z
    Duration
    property, 434
    2024 - 07 - 0
    9
    T11: 14:45.581953216
    Z
    From, To, and Duration
    properties, 430
    2024 - 07 - 0
    9
    T11: 14:45.581954049
    Z
    IsCumulative
    property, 439
    2024 - 07 - 0
    9
    T11: 14:45.581954966
    Z
    naming
    format, 426
    2024 - 07 - 0
    9
    T11: 14:45.581955799
    Z
    omitting
    both
    From and To
    properties, 432
    2024 - 07 - 0
    9
    T11: 14:45.581956716
    Z
    similarity
    of
    Duration
    property and
    2024 - 07 - 0
    9
    T11: 14:45.581957632
    Z
    TimeSpan
    object, 434
    2024 - 07 - 0
    9
    T11: 14:45.581963341
    Z
    using
    BeginAnimation()
    method
    to
    launch
    2024 - 07 - 0
    9
    T11: 14:45.581964674
    Z
    more
    than
    one
    animation
    at
    time, 434
    2024 - 07 - 0
    9
    T11: 14:45.581965716
    Z
    using
    By
    property
    instead
    of
    To
    property,
    2024 - 07 - 0
    9
    T11: 14:45.581966632
    Z
    433
    2024 - 07 - 0
    9
    T11: 14:45.581967632
    Z
    linear
    key
    frames, naming
    format, 478
    2024 - 07 - 0
    9
    T11: 14:45.581968549
    Z
    LinearGradientBrush, 162, 352, 449, 474
    2024 - 07 - 0
    9
    T11: 14:45.581969424
    Z
    changing
    lighting or color, 521
    2024 - 07 - 0
    9
    T11: 14:45.581970341
    Z
    creating
    blended
    fill, 354
    2024 - 07 - 0
    9
    T11: 14:45.581971216
    Z
    markup
    for shading rectangle diagonally,
    2024 - 07 - 0
    9
    T11: 14:45.581972132
    Z
    354
    2024 - 07 - 0
    9
    T11: 14:45.581972966
    Z
    proportional
    coordinate
    system, 355
    2024 - 07 - 0
    9
    T11: 14:45.581973841
    Z
    SpreadMethod
    property, 355
    2024 - 07 - 0
    9
    T11: 14:45.581975132
    Z
    using
    StartPoint and EndPoint
    properties,
    2024 - 07 - 0
    9
    T11: 14:45.581976132
    Z
    355
    2024 - 07 - 0
    9
    T11: 14:45.581977007
    Z
    LineBreak
    element, 950
    2024 - 07 - 0
    9
    T11: 14:45.581977924
    Z
    LineBreakBefore
    property,
    2024 - 07 - 0
    9
    T11: 14:45.581981257
    Z
    FrameworkPropertyMetadata
    object, 587
    2024 - 07 - 0
    9
    T11: 14:45.581982424
    Z
    LineCount
    property, 198
    2024 - 07 - 0
    9
    T11: 14:45.581983299
    Z
    LineDown()
    method, 190
    2024 - 07 - 0
    9
    T11: 14:45.581984174
    Z
    LineGeometry

    class , 375, 377

    2024 - 07 - 0
    9
    T11: 14:45.581985174
    Z
    LineHeight
    property, 940, 1009
    2024 - 07 - 0
    9
    T11: 14:45.581986382
    Z
    LineLeft()
    method, 190
    2024 - 07 - 0
    9
    T11: 14:45.581987257
    Z
    LineRight()
    method, 190
    2024 - 07 - 0
    9
    T11: 14:45.581988132
    Z
    LineSegment

    class , 384

    2024 - 07 - 0
    9
    T11: 14:45.581989007
    Z
    LineStackingStrategy
    property, 940
    2024 - 07 - 0
    9
    T11: 14:45.581989966
    Z
    LineUp()
    method, 190
    2024 - 07 - 0
    9
    T11: 14:45.581995882
    Z
    2024 - 07 - 0
    9
    T11: 14:45.581997924
    Z
    2024 - 07 - 0
    9
    T11: 14:45.581998924
    Z - ----
    2024 - 07 - 0
    9
    T11: 14:45.581999799
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582000632
    Z
    LinkLabel
    control, 1020
    2024 - 07 - 0
    9
    T11: 14:45.582001507
    Z
    LINQ(Language
    Integrated
    Query), 617—618
    2024 - 07 - 0
    9
    T11: 14:45.582002507
    Z
    list
    controls
    2024 - 07 - 0
    9
    T11: 14:45.582003382
    Z
    ComboBox
    control, 206
    2024 - 07 - 0
    9
    T11: 14:45.582004216
    Z
    ItemsControl

    class , 202

    2024 - 07 - 0
    9
    T11: 14:45.582005174
    Z
    ListBox
    control, 203
    2024 - 07 - 0
    9
    T11: 14:45.582006132
    Z
    overview, 159
    2024 - 07 - 0
    9
    T11: 14:45.582006966
    Z
    List
    element, 944
    2024 - 07 - 0
    9
    T11: 14:45.582007841
    Z
    ListBox

    class , 203

    2024 - 07 - 0
    9
    T11: 14:45.582008716
    Z
    ListBox
    control, 1068
    2024 - 07 - 0
    9
    T11: 14:45.582009591
    Z
    binding
    expression
    for
        2024 - 07 - 0
        9
        T11: 14:45.582010424
        Z
        RadioButton.IsChecked
        property, 662
    2024 - 07 - 0
    9
    T11: 14:45.582011341
    Z
    changing
    SelectionMode
    property
    to
    allow
    2024 - 07 - 0
    9
    T11: 14:45.582012299
    Z
    multiple
    selection, 663
    2024 - 07 - 0
    9
    T11: 14:45.582013341
    Z
    CheckBox
    element, 660, 663
    2024 - 07 - 0
    9
    T11: 14:45.582014257
    Z
    combining
    text and image
    content in, 204
    2024 - 07 - 0
    9
    T11: 14:45.582015174
    Z
    ContainerFromElement()
    method, 206
    2024 - 07 - 0
    9
    T11: 14:45.582016007
    Z
    ContentPresenter
    element, 662
    2024 - 07 - 0
    9
    T11: 14:45.582016924
    Z
    displaying
    check
    boxes in, 663
    2024 - 07 - 0
    9
    T11: 14:45.582017799
    Z
    DisplayMemberPath
    property, 662
    2024 - 07 - 0
    9
    T11: 14:45.582018674
    Z
    IsSelected
    property, 206
    2024 - 07 - 0
    9
    T11: 14:45.582019549
    Z
    ItemContainerStyle
    property, 663
    2024 - 07 - 0
    9
    T11: 14:45.582020466
    Z
    Items
    collection, 203
    2024 - 07 - 0
    9
    T11: 14:45.582021341
    Z
    ItemTemplate
    property, 662
    2024 - 07 - 0
    9
    T11: 14:45.582022257
    Z
    ListBoxItem.Control
    template, 663
    2024 - 07 - 0
    9
    T11: 14:45.582023132
    Z
    manually
    placing
    items in list, 206
    2024 - 07 - 0
    9
    T11: 14:45.582024049
    Z
    modifying
    ListBoxItem.Template
    property,
    2024 - 07 - 0
    9
    T11: 14:45.582024966
    Z
    662
    2024 - 07 - 0
    9
    T11: 14:45.582025757
    Z
    modifyingListBoxItem.Template
    property,
    2024 - 07 - 0
    9
    T11: 14:45.582026674
    Z
    662
    2024 - 07 - 0
    9
    T11: 14:45.582027549
    Z
    nesting
    arbitrary
    elements
    inside
    list
    box
    2024 - 07 - 0
    9
    T11: 14:45.582028466
    Z
    items, 204
    2024 - 07 - 0
    9
    T11: 14:45.582029341
    Z
    RadioButton
    element, 660, 663
    2024 - 07 - 0
    9
    T11: 14:45.582030216
    Z
    RemovedItems
    property, 205
    2024 - 07 - 0
    9
    T11: 14:45.582031716
    Z
    retrieving
    ListBoxItem
    wrapper
    for specific
        2024 - 07 - 0
        9
        T11: 14:45.582032632
        Z
        object, 206
    2024 - 07 - 0
    9
    T11: 14:45.582033466
    Z
    Selected
    event, 206
    2024 - 07 - 0
    9
    T11: 14:45.582034341
    Z
    SelectedItem
    property, 205
    2024 - 07 - 0
    9
    T11: 14:45.582035216
    Z
    SelectedItems
    property, 658
    2024 - 07 - 0
    9
    T11: 14:45.582036466
    Z
    SelectionChanged
    event, 205—206
    2024 - 07 - 0
    9
    T11: 14:45.582039049
    Z
    SelectionMode
    property, 658
    2024 - 07 - 0
    9
    T11: 14:45.582039882
    Z
    setting
    RadioButton.Focusable
    property, 662
    2024 - 07 - 0
    9
    T11: 14:45.582040799
    Z
    Unselected
    event, 206
    2024 - 07 - 0
    9
    T11: 14:45.582041716
    Z
    ListBoxChrome

    class , 509

    2024 - 07 - 0
    9
    T11: 14:45.582042591
    Z
    ListBoxChrome
    decorator, 75
    2024 - 07 - 0
    9
    T11: 14:45.582043424
    Z
    ListBoxItem
    elements, 203
    2024 - 07 - 0
    9
    T11: 14:45.582044341
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582045174
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582046007
    Z
    ListCollectionView, 691, 693, 702
    2024 - 07 - 0
    9
    T11: 14:45.582046841
    Z
    ListSelectionJournalEntry
    callback, 813
    2024 - 07 - 0
    9
    T11: 14:45.582047716
    Z
    ListView

    class
        2024 - 07 - 0
        9
        T11: 14:45.582048549
        Z
        ControlTemplate, 715

    2024 - 07 - 0
    9
    T11: 14:45.582049424
    Z
    creating
    custom
    view, 715
    2024 - 07 - 0
    9
    T11: 14:45.582050299
    Z
    creating
    customizable
    multicolumned
    lists,
    2024 - 07 - 0
    9
    T11: 14:45.582054341
    Z
    710
    2024 - 07 - 0
    9
    T11: 14:45.582055299
    Z
    creating
    grid
    that
    can
    switch
    views, 716
    2024 - 07 - 0
    9
    T11: 14:45.582056257
    Z
    DataTemplate, 715
    2024 - 07 - 0
    9
    T11: 14:45.582057257
    Z
    DefaultStyleKey
    property, 715
    2024 - 07 - 0
    9
    T11: 14:45.582058132
    Z
    function
    of, 710
    2024 - 07 - 0
    9
    T11: 14:45.582059007
    Z
    ItemContainerDefaultKeyStyle
    property, 715
    2024 - 07 - 0
    9
    T11: 14:45.582059882
    Z
    ResourceKey
    object, 715
    2024 - 07 - 0
    9
    T11: 14:45.582060716
    Z
    separating
    ListView
    control
    from View
    2024 - 07 - 0
    9
    T11: 14:45.582061757
    Z
    objects, 710
    2024 - 07 - 0
    9
    T11: 14:45.582062591
    Z
    switching
    between
    multiple
    views
    with same
        2024 - 07 - 0
        9
        T11: 14:45.582063466
        Z
        list, 710
    2024 - 07 - 0
    9
    T11: 14:45.582064382
    Z
    TileView

    class , 717

    2024 - 07 - 0
    9
    T11: 14:45.582067757
    Z
    View
    property, 710
    2024 - 07 - 0
    9
    T11: 14:45.582068924
    Z
    View
    property, advantages
    of, 710
    2024 - 07 - 0
    9
    T11: 14:45.582069841
    Z
    ViewBase

    class , 710

    2024 - 07 - 0
    9
    T11: 14:45.582070716
    Z
    ListView
    control
    2024 - 07 - 0
    9
    T11: 14:45.582072132
    Z
    adding
    properties
    to
    view
    classes, 722
    2024 - 07 - 0
    9
    T11: 14:45.582073049
    Z
    adding
    Setter
    to
    replace
    ControlTemplate,
    2024 - 07 - 0
    9
    T11: 14:45.582073966
    Z
    723
    2024 - 07 - 0
    9
    T11: 14:45.582074799
    Z
    defining
    view
    objects in Windows.Resources
    2024 - 07 - 0
    9
    T11: 14:45.582075716
    Z
    collection, 721
    2024 - 07 - 0
    9
    T11: 14:45.582076674
    Z
    GridView

    class , 721

    2024 - 07 - 0
    9
    T11: 14:45.582077591
    Z
    ImageDetailView
    object, 721
    2024 - 07 - 0
    9
    T11: 14:45.582078466
    Z
    ImageView
    object, 721
    2024 - 07 - 0
    9
    T11: 14:45.582079466
    Z
    passing
    information
    to
    view, 722
    2024 - 07 - 0
    9
    T11: 14:45.582080341
    Z
    setting
    ListView.View
    property, 720
    2024 - 07 - 0
    9
    T11: 14:45.582081257
    Z
    using
    custom
    view, 720
    2024 - 07 - 0
    9
    T11: 14:45.582082091
    Z
    Load()
    method, 866
    2024 - 07 - 0
    9
    T11: 14:45.582084382
    Z
    LoadAsync()
    method, 866
    2024 - 07 - 0
    9
    T11: 14:45.582085216
    Z
    LoadCompleted
    event, 809, 866
    2024 - 07 - 0
    9
    T11: 14:45.582086091
    Z
    LoadComponent()
    method, 31
    2024 - 07 - 0
    9
    T11: 14:45.582086966
    Z
    Loaded
    event, 134, 776
    2024 - 07 - 0
    9
    T11: 14:45.582087841
    Z
    LoadedBehavior
    property, 871
    2024 - 07 - 0
    9
    T11: 14:45.582089007
    Z
    LoadFile()
    method, 224
    2024 - 07 - 0
    9
    T11: 14:45.582089882
    Z
    LoadingRow
    event, DataGrid, 742—743
    2024 - 07 - 0
    9
    T11: 14:45.582090799
    Z
    LocalizabilityAttribute, 243
    2024 - 07 - 0
    9
    T11: 14:45.582091674
    Z
    localization
    2024 - 07 - 0
    9
    T11: 14:45.582092549
    Z
    adding < PropertyGroup > element.csproj
    file,
    2024 - 07 - 0
    9
    T11: 14:45.582093674
    Z
    241
    2024 - 07 - 0
    9
    T11: 14:45.582094507
    Z
    adding
    specialized
    Uid
    attribute
    to
    2024 - 07 - 0
    9
    T11: 14:45.582095424
    Z
    elements, 242
    2024 - 07 - 0
    9
    T11: 14:45.582096299
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582097132
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582097924
    Z - ----
    2024 - 07 - 0
    9
    T11: 14:45.582098757
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582099591
    Z
    adding
    support
    for more than one culture to
    2024 - 07 - 0
    9
    T11: 14:45.582100466
    Z
    application, 242
    2024 - 07 - 0
    9
    T11: 14:45.582101299
    Z
    building
    localizable
    user
    interfaces, 240
    2024 - 07 - 0
    9
    T11: 14:45.582102216
    Z
    building
    satellite
    assembly, 246
    2024 - 07 - 0
    9
    T11: 14:45.582103132
    Z
    culture
    names and their
    two - part
    identifiers,
    2024 - 07 - 0
    9
    T11: 14:45.582104007
    Z
    241
    2024 - 07 - 0
    9
    T11: 14:45.582104841
    Z
    CultureInfo

    class , 246

    2024 - 07 - 0
    9
    T11: 14:45.582105716
    Z
    CurrentUICulture
    property, 240
    2024 - 07 - 0
    9
    T11: 14:45.582106591
    Z
    extracting
    localizable
    content, 244
    2024 - 07 - 0
    9
    T11: 14:45.582107466
    Z
    global assembly
    cache(GAC), 242
    2024 - 07 - 0
    9
    T11: 14:45.582108507
    Z
    Global
    Sans
    Serif, 241
    2024 - 07 - 0
    9
    T11: 14:45.582109424
    Z
    Global
    Serif, 241
    2024 - 07 - 0
    9
    T11: 14:45.582110299
    Z
    Global
    User
    Interface, 241
    2024 - 07 - 0
    9
    T11: 14:45.582111174
    Z
    LocalizabilityAttribute, 243
    2024 - 07 - 0
    9
    T11: 14:45.582112091
    Z
    localizing
    FontFamily
    property in user
    2024 - 07 - 0
    9
    T11: 14:45.582113007
    Z
    interface, 241
    2024 - 07 - 0
    9
    T11: 14:45.582114382
    Z
    managing
    localization
    process, 242
    2024 - 07 - 0
    9
    T11: 14:45.582115299
    Z
    placing
    localized
    BAML
    resources in satellite
    2024 - 07 - 0
    9
    T11: 14:45.582116257
    Z
    assemblies, 240
    2024 - 07 - 0
    9
    T11: 14:45.582117257
    Z
    preparing
    application
    for , 241
    2024 - 07 - 0
    9
    T11: 14:45.582118132
    Z
    preparing
    markup
    elements
    for , 242
    2024 - 07 - 0
    9
    T11: 14:45.582119049
    Z
    probing, 240
    2024 - 07 - 0
    9
    T11: 14:45.582122757
    Z
    setting
    FlowDirection
    property
    for right - toleft layouts, 241
        2024 - 07 - 0
        9
        T11: 14:45.582123757
        Z
        using
        LocBaml.exe
        command - line
        tool, 244
    2024 - 07 - 0
    9
    T11: 14:45.582124674
    Z
    using
    MSBuild
    to
    generate
    Uid
    attributes,
    2024 - 07 - 0
    9
    T11: 14:45.582125716
    Z
    243
    2024 - 07 - 0
    9
    T11: 14:45.582128549
    Z
    XAML
    file as unit
    of
    localization, 240
    2024 - 07 - 0
    9
    T11: 14:45.582129549
    Z
    localized
    text, 101
    2024 - 07 - 0
    9
    T11: 14:45.582130424
    Z
    LocalPrintServer

    class , 1013

    2024 - 07 - 0
    9
    T11: 14:45.582131341
    Z
    LocationChanged
    event, 754
    2024 - 07 - 0
    9
    T11: 14:45.582132257
    Z
    LocBaml.exe
    2024 - 07 - 0
    9
    T11: 14:45.582133091
    Z
    building
    satellite
    assembly, 246
    2024 - 07 - 0
    9
    T11: 14:45.582133966
    Z
    compiling
    by
    hand, 244
    2024 - 07 - 0
    9
    T11: 14:45.582134966
    Z / cul: parameter, 246
    2024 - 07 - 0
    9
    T11: 14:45.582135841
    Z / generate
    parameter, 246
    2024 - 07 - 0
    9
    T11: 14:45.582136757
    Z / parse
    parameter, 244
    2024 - 07 - 0
    9
    T11: 14:45.582137841
    Z
    table
    of
    localizable
    properties, 245
    2024 - 07 - 0
    9
    T11: 14:45.582138757
    Z / trans: parameter, 246
    2024 - 07 - 0
    9
    T11: 14:45.582139632
    Z
    Lock()
    method, 421
    2024 - 07 - 0
    9
    T11: 14:45.582140549
    Z
    logical
    scrolling, 191
    2024 - 07 - 0
    9
    T11: 14:45.582141424
    Z
    logical
    tree
    2024 - 07 - 0
    9
    T11: 14:45.582142257
    Z
    building, 500
    2024 - 07 - 0
    9
    T11: 14:45.582143132
    Z
    LogicalTreeHelper

    class , table of methods,

    2024 - 07 - 0
    9
    T11: 14:45.582144091
    Z
    503
    2024 - 07 - 0
    9
    T11: 14:45.582144924
    Z
    LogicalTreeHelper, 52—53, 503
    2024 - 07 - 0
    9
    T11: 14:45.582145841
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582146674
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582147507
    Z
    Long
    Date
    data
    types, data
    binding
    format
    2024 - 07 - 0
    9
    T11: 14:45.582148382
    Z
    string, 647
    2024 - 07 - 0
    9
    T11: 14:45.582149257
    Z
    lookless
    controls, 21
    2024 - 07 - 0
    9
    T11: 14:45.582150132
    Z
    adding
    TemplatePart
    attribute
    to
    control
    2024 - 07 - 0
    9
    T11: 14:45.582155966
    Z
    declaration, 566
    2024 - 07 - 0
    9
    T11: 14:45.582156882
    Z
    calling
    OverrideMetadata()
    method, 560
    2024 - 07 - 0
    9
    T11: 14:45.582157799
    Z
    changing
    color
    picker
    into
    lookless
    control,
    2024 - 07 - 0
    9
    T11: 14:45.582158674
    Z
    560
    2024 - 07 - 0
    9
    T11: 14:45.582159632
    Z
    checking
    for correct type of element, 564
    2024 - 07 - 0
    9
    T11: 14:45.582160674
    Z
    code
    for binding SolidColorBrush, 565
        2024 - 07 - 0
        9
        T11: 14:45.582161507
        Z
        connecting
        data
        binding
        expression
        using
    2024 - 07 - 0
    9
    T11: 14:45.582162424
    Z
    OnApplyTemplate()
    method, 564
    2024 - 07 - 0
    9
    T11: 14:45.582163424
    Z
    converting
    ordinary
    markup
    into
    control
    2024 - 07 - 0
    9
    T11: 14:45.582164299
    Z
    template, 562
    2024 - 07 - 0
    9
    T11: 14:45.582165132
    Z
    creating, 559
    2024 - 07 - 0
    9
    T11: 14:45.582165966
    Z
    creating
    template
    for color picker, 562
        2024 - 07 - 0
        9
        T11: 14:45.582166882
        Z
        DefaultStyleKeyProperty, 560
    2024 - 07 - 0
    9
    T11: 14:45.582167757
    Z
    definition
    of, 506, 559
    2024 - 07 - 0
    9
    T11: 14:45.582168632
    Z
    ElementName
    property, 562
    2024 - 07 - 0
    9
    T11: 14:45.582169507
    Z
    generic.xaml
    resource
    dictionary, 560
    2024 - 07 - 0
    9
    T11: 14:45.582170632
    Z
    markup
    structure
    for ColorPicker.xaml, 561
        2024 - 07 - 0
        9
        T11: 14:45.582173507
        Z
        providing
        descriptive
        names
        for element
            2024 - 07 - 0
        9
        T11: 14:45.582174382
        Z
        names, 563
    2024 - 07 - 0
    9
    T11: 14:45.582175216
    Z
    RelativeSource
    property, 562
    2024 - 07 - 0
    9
    T11: 14:45.582176049
    Z
    streamlining
    color
    picker
    control
    template,
    2024 - 07 - 0
    9
    T11: 14:45.582176966
    Z
    563
    2024 - 07 - 0
    9
    T11: 14:45.582177799
    Z
    TemplateBinding, 562—563
    2024 - 07 - 0
    9
    T11: 14:45.582178674
    Z
    using
    TargetType
    attribute, 561
    2024 - 07 - 0
    9
    T11: 14:45.582179632
    Z
    lookless
    WPF
    controls, 27
    2024 - 07 - 0
    9
    T11: 14:45.582180507
    Z
    loose
    XAML
    files, opening in Internet
    Explorer,
    2024 - 07 - 0
    9
    T11: 14:45.582181424
    Z
    55
    2024 - 07 - 0
    9
    T11: 14:45.582182257
    Z
    LostFocus
    event, 135
    2024 - 07 - 0
    9
    T11: 14:45.582183216
    Z
    LostFocus
    update
    mode, 259
    2024 - 07 - 0
    9
    T11: 14:45.582184591
    Z
    LostMouseCapture
    event, 146
    2024 - 07 - 0
    9
    T11: 14:45.582185466
    Z
    LowerLatin
    value, MarkerStyle
    property, 944
    2024 - 07 - 0
    9
    T11: 14:45.582186382
    Z
    LowerRoman
    value, MarkerStyle
    property, 944
    2024 - 07 - 0
    9
    T11: 14:45.582190382
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582192007
    Z ■ ** M **
    2024 - 07 - 0
    9
    T11: 14:45.582193007
    Z
    Mad
    Libs
    game, creating, 957
    2024 - 07 - 0
    9
    T11: 14:45.582193882
    Z
    MAF(Managed
    Add - in Framework), 1055—1056
    2024 - 07 - 0
    9
    T11: 14:45.582194882
    Z
    Main()
    method, 216—217, 220
    2024 - 07 - 0
    9
    T11: 14:45.582195757
    Z
    main
    page, bomb - dropping
    game, 488
    2024 - 07 - 0
    9
    T11: 14:45.582197799
    Z
    MainWindow
    property, 216, 758
    2024 - 07 - 0
    9
    T11: 14:45.582200507
    Z
    Managed
    Add - in Framework(MAF), 1055—1056
    2024 - 07 - 0
    9
    T11: 14:45.582201632
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582202507
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582203341
    Z - ----
    2024 - 07 - 0
    9
    T11: 14:45.582204174
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582205007
    Z
    Managed
    Extensibility
    Framework(MEF),
    2024 - 07 - 0
    9
    T11: 14:45.582205882
    Z
    1055—1056
    2024 - 07 - 0
    9
    T11: 14:45.582206799
    Z
    manifests, 233
    2024 - 07 - 0
    9
    T11: 14:45.582207632
    Z
    manipulation, 153—156
    2024 - 07 - 0
    9
    T11: 14:45.582208841
    Z
    Margin
    property, 68, 940
    2024 - 07 - 0
    9
    T11: 14:45.582209757
    Z
    margins
    2024 - 07 - 0
    9
    T11: 14:45.582210674
    Z
    keeping
    settings
    consistent, 70
    2024 - 07 - 0
    9
    T11: 14:45.582211549
    Z
    setting
    for StackPanel, 69
        2024 - 07 - 0
        9
        T11: 14:45.582212424
        Z
        Thickness
        structure, 70
    2024 - 07 - 0
    9
    T11: 14:45.582213299
    Z
    MarkerStyle
    property, 944
    2024 - 07 - 0
    9
    T11: 14:45.582214216
    Z
    markup
    extensions, using in nested
    tags or XML
    2024 - 07 - 0
    9
    T11: 14:45.582215132
    Z
    attributes, 37
    2024 - 07 - 0
    9
    T11: 14:45.582216049
    Z
    MarshalByRefObject
    attribute, 1066
    2024 - 07 - 0
    9
    T11: 14:45.582217049
    Z
    MaskedTextBox
    control, 1020, 1022
    2024 - 07 - 0
    9
    T11: 14:45.582217966
    Z
    MaskedTextBox, ValidatingType
    property, 1032
    2024 - 07 - 0
    9
    T11: 14:45.582218882
    Z
    Material

    class , 895

    2024 - 07 - 0
    9
    T11: 14:45.582222591
    Z
    Material
    property, 895
    2024 - 07 - 0
    9
    T11: 14:45.582223424
    Z
    MaterialGroup

    class , 895

    2024 - 07 - 0
    9
    T11: 14:45.582224424
    Z
    MatrixCamera, transforming
    3 - D
    scene
    to
    2 - D
    2024 - 07 - 0
    9
    T11: 14:45.582225341
    Z
    view, 899
    2024 - 07 - 0
    9
    T11: 14:45.582226341
    Z
    MatrixTransform

    class , 366

    2024 - 07 - 0
    9
    T11: 14:45.582227216
    Z
    MaxHeight
    property, 68
    2024 - 07 - 0
    9
    T11: 14:45.582228174
    Z
    MaxLength
    property, 198, 202
    2024 - 07 - 0
    9
    T11: 14:45.582229216
    Z
    MaxLines
    property, 198
    2024 - 07 - 0
    9
    T11: 14:45.582230091
    Z
    MaxWidth
    property, 68
    2024 - 07 - 0
    9
    T11: 14:45.582230966
    Z
    Maya, 910
    2024 - 07 - 0
    9
    T11: 14:45.582231841
    Z
    MDI(multiple
    document
    interface), 761, 1021—
    2024 - 07 - 0
    9
    T11: 14:45.582232757
    Z
    1022
    2024 - 07 - 0
    9
    T11: 14:45.582233716
    Z
    Measure()
    method, 583, 994
    2024 - 07 - 0
    9
    T11: 14:45.582234716
    Z
    measure
    pass, 583
    2024 - 07 - 0
    9
    T11: 14:45.582235591
    Z
    measure
    stage, 63
    2024 - 07 - 0
    9
    T11: 14:45.582236466
    Z
    MeasureCore()
    method, 583
    2024 - 07 - 0
    9
    T11: 14:45.582238216
    Z
    MeasureOverride()
    method, 64, 586
    2024 - 07 - 0
    9
    T11: 14:45.582239091
    Z
    allowing
    child
    to
    take
    all
    space
    it
    wants, 584
    2024 - 07 - 0
    9
    T11: 14:45.582240049
    Z
    basic
    structure
    of, 583
    2024 - 07 - 0
    9
    T11: 14:45.582240924
    Z
    calling
    Measure()
    method
    of
    each
    child, 583
    2024 - 07 - 0
    9
    T11: 14:45.582241966
    Z
    DesiredSize
    property, 584
    2024 - 07 - 0
    9
    T11: 14:45.582243132
    Z
    determining
    how
    much
    space
    each
    child
    2024 - 07 - 0
    9
    T11: 14:45.582244132
    Z
    wants, 583
    2024 - 07 - 0
    9
    T11: 14:45.582245049
    Z
    passing
    Size
    object
    with value of
    2024 - 07 - 0
    9
    T11: 14:45.582245966
    Z
    Double.PositiveInfinity, 584
    2024 - 07 - 0
    9
    T11: 14:45.582246924
    Z
    Media
    Integration
    Layer(MIL), 13
    2024 - 07 - 0
    9
    T11: 14:45.582247841
    Z
    MediaClock

    class , 873

    2024 - 07 - 0
    9
    T11: 14:45.582248716
    Z
    MediaCommands

    class , types of included

    2024 - 07 - 0
    9
    T11: 14:45.582249799
    Z
    commands, 270
    2024 - 07 - 0
    9
    T11: 14:45.582250632
    Z
    MediaElement

    class , 239, 1025

    2024 - 07 - 0
    9
    T11: 14:45.582251716
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582252549
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582253341
    Z
    adding
    MediaElement
    tag
    for playing sound,
        2024 - 07 - 0
        9
        T11: 14:45.582254257
        Z
        871
    2024 - 07 - 0
    9
    T11: 14:45.582255091
    Z
    Balance
    property, 877
    2024 - 07 - 0
    9
    T11: 14:45.582256007
    Z
    Clock
    property, 873
    2024 - 07 - 0
    9
    T11: 14:45.582257091
    Z
    controlling
    additional
    playback
    details, 876
    2024 - 07 - 0
    9
    T11: 14:45.582260757
    Z
    controlling
    audio
    declaratively
    through
    2024 - 07 - 0
    9
    T11: 14:45.582261841
    Z
    XAML
    markup, 872
    2024 - 07 - 0
    9
    T11: 14:45.582262799
    Z
    controlling
    audio
    playback
    2024 - 07 - 0
    9
    T11: 14:45.582263674
    Z
    programmatically, 871
    2024 - 07 - 0
    9
    T11: 14:45.582264591
    Z
    creating
    video - reflection
    effect, code
    2024 - 07 - 0
    9
    T11: 14:45.582265507
    Z
    example, 881
    2024 - 07 - 0
    9
    T11: 14:45.582267966
    Z
    error
    handling, 872
    2024 - 07 - 0
    9
    T11: 14:45.582274049
    Z
    ErrorException
    property, 872
    2024 - 07 - 0
    9
    T11: 14:45.582275257
    Z
    LayoutTransform
    property, 881
    2024 - 07 - 0
    9
    T11: 14:45.582276216
    Z
    LoadedBehavior
    property, 871
    2024 - 07 - 0
    9
    T11: 14:45.582277091
    Z
    Manual
    mode, 878
    2024 - 07 - 0
    9
    T11: 14:45.582277924
    Z
    MediaState
    enumeration, 871
    2024 - 07 - 0
    9
    T11: 14:45.582278924
    Z
    Pause()
    method, 872
    2024 - 07 - 0
    9
    T11: 14:45.582279966
    Z
    placement
    of,
    for audio and video, 871
    2024 - 07 - 0
    9
    T11: 14:45.582280882
    Z
    Play()
    method, 872
    2024 - 07 - 0
    9
    T11: 14:45.582281841
    Z
    playing
    audio
    with triggers, 872
        2024 - 07 - 0
        9
        T11: 14:45.582282757
        Z
        playing
        multiple
        audio
        files, code
        example,
    2024 - 07 - 0
    9
    T11: 14:45.582283716
    Z
    875
    2024 - 07 - 0
    9
    T11: 14:45.582284591
    Z
    playing
    video, 880
    2024 - 07 - 0
    9
    T11: 14:45.582285424
    Z
    Position
    property, 878
    2024 - 07 - 0
    9
    T11: 14:45.582286299
    Z
    RenderTransform
    property, 881
    2024 - 07 - 0
    9
    T11: 14:45.582287216
    Z
    RenderTransformOrigin
    property, 881
    2024 - 07 - 0
    9
    T11: 14:45.582288091
    Z
    setting
    Clipping
    property, 881
    2024 - 07 - 0
    9
    T11: 14:45.582289174
    Z
    setting
    Opacity
    property, 881
    2024 - 07 - 0
    9
    T11: 14:45.582290507
    Z
    setting
    Position
    to
    move
    through
    audio
    file,
    2024 - 07 - 0
    9
    T11: 14:45.582291466
    Z
    872
    2024 - 07 - 0
    9
    T11: 14:45.582292341
    Z
    SpeedRatio
    property, 877
    2024 - 07 - 0
    9
    T11: 14:45.582293216
    Z
    Stop()
    method, 872
    2024 - 07 - 0
    9
    T11: 14:45.582294091
    Z
    Stretch
    property, 880
    2024 - 07 - 0
    9
    T11: 14:45.582295632
    Z
    StretchDirection
    property, 880
    2024 - 07 - 0
    9
    T11: 14:45.582296591
    Z
    support
    for WMV, MPEG, and AVI files, 880
    2024 - 07 - 0
    9
    T11: 14:45.582297632
    Z
    synchronizing
    animation
    with audio or
        2024 - 07 - 0
        9
        T11: 14:45.582298799
        Z
        video
        file, 878
    2024 - 07 - 0
    9
    T11: 14:45.582299632
    Z
    types
    of
    video
    effects, 881
    2024 - 07 - 0
    9
    T11: 14:45.582300507
    Z
    using
    separate
    ResumeStoryboard
    action
    2024 - 07 - 0
    9
    T11: 14:45.582301466
    Z
    after
    pausing
    playback, 875
    2024 - 07 - 0
    9
    T11: 14:45.582302341
    Z
    using
    single
    Storyboard
    to
    control
    audio
    2024 - 07 - 0
    9
    T11: 14:45.582303299
    Z
    playback, 873
    2024 - 07 - 0
    9
    T11: 14:45.582306757
    Z
    Volume
    property, 877
    2024 - 07 - 0
    9
    T11: 14:45.582308007
    Z
    MediaFailed
    event, 869
    2024 - 07 - 0
    9
    T11: 14:45.582308882
    Z
    MediaOpened
    event, 869
    2024 - 07 - 0
    9
    T11: 14:45.582309757
    Z
    MediaPlayer

    class , 1025

    2024 - 07 - 0
    9
    T11: 14:45.582310632
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582311591
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582312424
    Z - ----
    2024 - 07 - 0
    9
    T11: 14:45.582313299
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582314091
    Z
    creating
    Window.Unloaded
    event
    handler
    to
    2024 - 07 - 0
    9
    T11: 14:45.582315216
    Z
    call
    Close(), 869
    2024 - 07 - 0
    9
    T11: 14:45.582316091
    Z
    lack
    of
    exception
    handling
    code, 869
    2024 - 07 - 0
    9
    T11: 14:45.582318674
    Z
    MediaFailed
    event, 869
    2024 - 07 - 0
    9
    T11: 14:45.582319591
    Z
    MediaOpened
    event, 869
    2024 - 07 - 0
    9
    T11: 14:45.582320507
    Z
    no
    option
    for synchronous playback, 869
        2024 - 07 - 0
        9
        T11: 14:45.582321382
        Z
        Open()
        method, 869
    2024 - 07 - 0
    9
    T11: 14:45.582322382
    Z
    Play()
    method, 869
    2024 - 07 - 0
    9
    T11: 14:45.582323216
    Z
    playing
    multiple
    audio
    files, 876
    2024 - 07 - 0
    9
    T11: 14:45.582325299
    Z
    supplying
    location
    of
    audio
    file as URI, 869
    2024 - 07 - 0
    9
    T11: 14:45.582326174
    Z
    table
    of
    useful
    methods, properties, and
    2024 - 07 - 0
    9
    T11: 14:45.582327091
    Z
    events, 869
    2024 - 07 - 0
    9
    T11: 14:45.582327924
    Z
    MediaState
    enumeration, 871
    2024 - 07 - 0
    9
    T11: 14:45.582328841
    Z
    MediaTimeline

    class , 437, 873

    2024 - 07 - 0
    9
    T11: 14:45.582329716
    Z
    MEF(Managed
    Extensibility
    Framework),
    2024 - 07 - 0
    9
    T11: 14:45.582330632
    Z
    1055—1056
    2024 - 07 - 0
    9
    T11: 14:45.582331507
    Z
    MemoryStream, 978
    2024 - 07 - 0
    9
    T11: 14:45.582332466
    Z
    Menu

    class
        2024 - 07 - 0
        9
        T11: 14:45.582333341
        Z
        creating
        scrollable
        sidebar
        menu, 842

    2024 - 07 - 0
    9
    T11: 14:45.582334216
    Z
    DisplayMemberPath
    property, 842
    2024 - 07 - 0
    9
    T11: 14:45.582335091
    Z
    dividing
    menus
    into
    groups
    of
    related
    2024 - 07 - 0
    9
    T11: 14:45.582336007
    Z
    commands, 846
    2024 - 07 - 0
    9
    T11: 14:45.582336924
    Z
    example
    of
    Separator
    that
    defines
    text
    title,
    2024 - 07 - 0
    9
    T11: 14:45.582337924
    Z
    846
    2024 - 07 - 0
    9
    T11: 14:45.582338841
    Z
    IsMainMenu
    property, 842
    2024 - 07 - 0
    9
    T11: 14:45.582339799
    Z
    ItemsSource
    property, 842
    2024 - 07 - 0
    9
    T11: 14:45.582340716
    Z
    ItemTemplate
    property, 842
    2024 - 07 - 0
    9
    T11: 14:45.582341591
    Z
    ItemTemplateSelector
    property, 842
    2024 - 07 - 0
    9
    T11: 14:45.582342466
    Z
    Separator as not content
    control, 847
    2024 - 07 - 0
    9
    T11: 14:45.582343466
    Z
    using
    menu
    separators, 846
    2024 - 07 - 0
    9
    T11: 14:45.582344466
    Z
    MenuItem

    class
        2024 - 07 - 0
        9
        T11: 14:45.582345382
        Z
        Command
        property, 844

    2024 - 07 - 0
    9
    T11: 14:45.582346216
    Z
    CommandParameter
    property, 844
    2024 - 07 - 0
    9
    T11: 14:45.582347216
    Z
    CommandTarget
    property, 844
    2024 - 07 - 0
    9
    T11: 14:45.582348466
    Z
    creating
    rudimentary
    menu
    structure, 843
    2024 - 07 - 0
    9
    T11: 14:45.582349382
    Z
    displaying
    check
    mark
    next
    to
    menu
    item,
    2024 - 07 - 0
    9
    T11: 14:45.582350299
    Z
    845
    2024 - 07 - 0
    9
    T11: 14:45.582351132
    Z
    handling
    Click
    event, 844
    2024 - 07 - 0
    9
    T11: 14:45.582352049
    Z
    having
    non - MenuItem
    objects
    inside
    Menu
    2024 - 07 - 0
    9
    T11: 14:45.582352924
    Z or MenuItem, 844
    2024 - 07 - 0
    9
    T11: 14:45.582353799
    Z
    Icon
    property, 845
    2024 - 07 - 0
    9
    T11: 14:45.582354757
    Z
    including
    keyboard
    shortcuts, 843
    2024 - 07 - 0
    9
    T11: 14:45.582355632
    Z
    InputGestureText
    property, 845
    2024 - 07 - 0
    9
    T11: 14:45.582356507
    Z
    IsChecked
    property, 845
    2024 - 07 - 0
    9
    T11: 14:45.582357341
    Z
    Separator
    objects, 843
    2024 - 07 - 0
    9
    T11: 14:45.582359382
    Z
    setting
    shortcut
    text
    for menu item, 845
        2024 - 07 - 0
        9
        T11: 14:45.582360299
        Z
        showing
        thumbnail
        icon, 845
    2024 - 07 - 0
    9
    T11: 14:45.582361174
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582362007
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582362882
    Z
    StaysOpenOnClick
    property, 844
    2024 - 07 - 0
    9
    T11: 14:45.582363757
    Z
    MergedDictionaries
    collection, 522, 528
    2024 - 07 - 0
    9
    T11: 14:45.582364632
    Z
    MergedDictionaries
    property, 303
    2024 - 07 - 0
    9
    T11: 14:45.582365507
    Z
    mesh, building
    basic, 893
    2024 - 07 - 0
    9
    T11: 14:45.582366382
    Z
    MeshGeometry

    class , 914

    2024 - 07 - 0
    9
    T11: 14:45.582367257
    Z
    MeshGeometry3D

    class
        2024 - 07 - 0
        9
        T11: 14:45.582368132
        Z
        Normals
        property, 893—894

    2024 - 07 - 0
    9
    T11: 14:45.582369049
    Z
    Positions
    property, 893—894
    2024 - 07 - 0
    9
    T11: 14:45.582370341
    Z
    table
    of
    properties, 893
    2024 - 07 - 0
    9
    T11: 14:45.582371174
    Z
    TextureCoordinates
    property, 893—895
    2024 - 07 - 0
    9
    T11: 14:45.582372091
    Z
    TriangleIndices
    property, 893—894
    2024 - 07 - 0
    9
    T11: 14:45.582372966
    Z
    MeshHit
    property, determining
    whether
    torus
    2024 - 07 - 0
    9
    T11: 14:45.582374007
    Z
    mesh
    has
    been
    hit, 926
    2024 - 07 - 0
    9
    T11: 14:45.582374924
    Z
    MessageBeep
    Win32
    API, 868
    2024 - 07 - 0
    9
    T11: 14:45.582375799
    Z
    MessageBox

    class , 762

    2024 - 07 - 0
    9
    T11: 14:45.582376674
    Z
    MessageBoxButton
    enumeration, 762
    2024 - 07 - 0
    9
    T11: 14:45.582377757
    Z
    MessageBoxImage
    enumeration, 762
    2024 - 07 - 0
    9
    T11: 14:45.582378799
    Z
    Microsoft
    Expression
    Blend, 23
    2024 - 07 - 0
    9
    T11: 14:45.582379674
    Z
    Microsoft
    HTML
    Object
    Library(mshtml.tlb),
    2024 - 07 - 0
    9
    T11: 14:45.582380591
    Z
    835
    2024 - 07 - 0
    9
    T11: 14:45.582381757
    Z
    Microsoft
    Installer(MSI), 1079—1080
    2024 - 07 - 0
    9
    T11: 14:45.582382674
    Z
    Microsoft
    Money, weblike
    interface
    of, 792
    2024 - 07 - 0
    9
    T11: 14:45.582383549
    Z
    Microsoft.NET
    2.0
    Framework
    Configuration
    2024 - 07 - 0
    9
    T11: 14:45.582384466
    Z
    Tool, 826
    2024 - 07 - 0
    9
    T11: 14:45.582385466
    Z
    Microsoft
    Office
    2007, creating
    XPS and PDF
    2024 - 07 - 0
    9
    T11: 14:45.582386382
    Z
    documents, 974
    2024 - 07 - 0
    9
    T11: 14:45.582387216
    Z
    Microsoft
    Speech
    Software
    Development
    Kit,
    2024 - 07 - 0
    9
    T11: 14:45.582388132
    Z
    887
    2024 - 07 - 0
    9
    T11: 14:45.582388966
    Z
    Microsoft
    Word, 228
    2024 - 07 - 0
    9
    T11: 14:45.582389841
    Z
    Microsoft, XPS(XML
    Paper
    Specification), 935,
    2024 - 07 - 0
    9
    T11: 14:45.582390716
    Z
    974
    2024 - 07 - 0
    9
    T11: 14:45.582391549
    Z
    Microsoft.Expression.Interactions.dll
    assembly
    2024 - 07 - 0
    9
    T11: 14:45.582392507
    Z
    design - time
    behavior
    support in Expression
    2024 - 07 - 0
    9
    T11: 14:45.582393632
    Z
    Blend, 330
    2024 - 07 - 0
    9
    T11: 14:45.582394466
    Z
    support
    for behaviors, 326
        2024 - 07 - 0
        9
        T11: 14:45.582395382
        Z
        Microsoft
        's Composite Application Library
    2024 - 07 - 0
    9
    T11: 14:45.582396257
    Z(CAL), 1056
    2024 - 07 - 0
    9
    T11: 14:45.582397132
    Z
    Microsoft.Win32
    namespace, 762, 1025
    2024 - 07 - 0
    9
    T11: 14:45.582399132
    Z
    Microsoft.Windows.Themes, 509
    2024 - 07 - 0
    9
    T11: 14:45.582400049
    Z
    MIL(Media
    Integration
    Layer), 13
    2024 - 07 - 0
    9
    T11: 14:45.582400924
    Z
    milcore.dll, 13
    2024 - 07 - 0
    9
    T11: 14:45.582401799
    Z
    MinHeight
    property, 68
    2024 - 07 - 0
    9
    T11: 14:45.582402674
    Z
    MinLines
    property, 198
    2024 - 07 - 0
    9
    T11: 14:45.582403507
    Z
    MinOrphanLines
    property, 965
    2024 - 07 - 0
    9
    T11: 14:45.582404382
    Z
    MinWidth
    property, 68
    2024 - 07 - 0
    9
    T11: 14:45.582405257
    Z
    MinWindowLines
    property, 965
    2024 - 07 - 0
    9
    T11: 14:45.582406091
    Z
    Miter
    line
    join, 349
    2024 - 07 - 0
    9
    T11: 14:45.582407091
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582409257
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582410216
    Z - ----
    2024 - 07 - 0
    9
    T11: 14:45.582411049
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582411841
    Z
    mnemonics, 160, 175, 1038
    2024 - 07 - 0
    9
    T11: 14:45.582412716
    Z
    modal
    windows, 754
    2024 - 07 - 0
    9
    T11: 14:45.582413591
    Z
    Mode
    property, 252
    2024 - 07 - 0
    9
    T11: 14:45.582414466
    Z
    Model3D

    class , 919

    2024 - 07 - 0
    9
    T11: 14:45.582415299
    Z
    Model3DGroup

    class , 910, 912

    2024 - 07 - 0
    9
    T11: 14:45.582416216
    Z
    modeless
    windows, 754
    2024 - 07 - 0
    9
    T11: 14:45.582417132
    Z
    ModelUIElement3D
    2024 - 07 - 0
    9
    T11: 14:45.582418007
    Z
    hit
    testing in, 928—929
    2024 - 07 - 0
    9
    T11: 14:45.582418924
    Z
    overview, 927—928
    2024 - 07 - 0
    9
    T11: 14:45.582419757
    Z
    ModelUIElement3D

    class , 925, 929

    2024 - 07 - 0
    9
    T11: 14:45.582420674
    Z
    ModelVisual3D

    class , 919, 925, 930

    2024 - 07 - 0
    9
    T11: 14:45.582421591
    Z
    modifier
    keys, checking
    status
    of, 142
    2024 - 07 - 0
    9
    T11: 14:45.582422466
    Z
    Mouse

    class , 145

    2024 - 07 - 0
    9
    T11: 14:45.582423466
    Z
    mouse
    cursors, 168—169
    2024 - 07 - 0
    9
    T11: 14:45.582424341
    Z
    mouse
    events, 133
    2024 - 07 - 0
    9
    T11: 14:45.582425216
    Z
    AllowDrop
    property, 148
    2024 - 07 - 0
    9
    T11: 14:45.582426049
    Z
    ButtonState
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582426924
    Z
    capturing
    mouse
    by
    calling
    Mouse.Capture(
        2024 - 07 - 0
    9
    T11: 14:45.582427841
    Z ), 146
    2024 - 07 - 0
    9
    T11: 14:45.582428716
    Z
    ClickCount
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582429591
    Z
    creating
    drag - and -drop
    source, 148
    2024 - 07 - 0
    9
    T11: 14:45.582430507
    Z
    direct
    events, definition
    of, 143
    2024 - 07 - 0
    9
    T11: 14:45.582431424
    Z
    drag - and -drop
    operations, 146
    2024 - 07 - 0
    9
    T11: 14:45.582432299
    Z
    DragDrop

    class , 147

    2024 - 07 - 0
    9
    T11: 14:45.582433174
    Z
    DragEnter
    event, 148
    2024 - 07 - 0
    9
    T11: 14:45.582434049
    Z
    dragging - and -dropping
    into
    other
    2024 - 07 - 0
    9
    T11: 14:45.582434924
    Z
    applications, 148
    2024 - 07 - 0
    9
    T11: 14:45.582435841
    Z
    getting
    mouse
    coordinates, 143
    2024 - 07 - 0
    9
    T11: 14:45.582436716
    Z
    IsMouseDirectlyOver
    property, 144
    2024 - 07 - 0
    9
    T11: 14:45.582437632
    Z
    IsMouseOver
    property, 144
    2024 - 07 - 0
    9
    T11: 14:45.582438799
    Z
    losing
    mouse
    capture, 146
    2024 - 07 - 0
    9
    T11: 14:45.582440799
    Z
    LostMouseCapture
    event, 146
    2024 - 07 - 0
    9
    T11: 14:45.582442424
    Z
    Mouse

    class , 145

    2024 - 07 - 0
    9
    T11: 14:45.582445466
    Z
    mouse
    click
    events
    for all elements, 144
        2024 - 07 - 0
        9
        T11: 14:45.582446382
        Z
        MouseButton
        event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582447257
    Z
    MouseButtonEventArgs
    object, 145
    2024 - 07 - 0
    9
    T11: 14:45.582448507
    Z
    MouseDoubleClick
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582449591
    Z
    MouseEnter
    event, 143
    2024 - 07 - 0
    9
    T11: 14:45.582450466
    Z
    MouseLeave
    event, 143
    2024 - 07 - 0
    9
    T11: 14:45.582451382
    Z
    MouseMove
    event, 143
    2024 - 07 - 0
    9
    T11: 14:45.582452257
    Z
    PreviewMouseDoubleClick
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582453174
    Z
    PreviewMouseMove
    event, 143
    2024 - 07 - 0
    9
    T11: 14:45.582454049
    Z
    state
    groups, 541
    2024 - 07 - 0
    9
    T11: 14:45.582454924
    Z
    MouseButton
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582455841
    Z
    MouseButtonEventArgs
    object, 127, 145, 928
    2024 - 07 - 0
    9
    T11: 14:45.582456757
    Z
    Mouse.Capture()
    method, 146
    2024 - 07 - 0
    9
    T11: 14:45.582457632
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582458507
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582459341
    Z
    MouseDoubleClick
    event, 145, 545
    2024 - 07 - 0
    9
    T11: 14:45.582460216
    Z
    MouseDown
    event, 124
    2024 - 07 - 0
    9
    T11: 14:45.582461091
    Z
    MouseEnter
    event, 124, 143, 316
    2024 - 07 - 0
    9
    T11: 14:45.582461966
    Z
    MouseEventArgs
    object, 121, 143
    2024 - 07 - 0
    9
    T11: 14:45.582462841
    Z
    MouseLeave
    event, 143, 316, 469
    2024 - 07 - 0
    9
    T11: 14:45.582463799
    Z
    MouseLeftButtonDown
    event, 145, 494, 768
    2024 - 07 - 0
    9
    T11: 14:45.582464716
    Z
    MouseLeftButtonUp
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582465674
    Z
    MouseMove
    event, 143
    2024 - 07 - 0
    9
    T11: 14:45.582466549
    Z
    MouseOver
    state, controls, 541
    2024 - 07 - 0
    9
    T11: 14:45.582467507
    Z
    MouseRightButtonDown
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582468382
    Z
    MouseRightButtonUp
    event, 145
    2024 - 07 - 0
    9
    T11: 14:45.582469257
    Z
    MouseUp
    event, 121
    2024 - 07 - 0
    9
    T11: 14:45.582470132
    Z
    MouseUp()
    method, 122
    2024 - 07 - 0
    9
    T11: 14:45.582471174
    Z
    MSBuild, using
    to
    generate
    Uid
    attributes, 243
    2024 - 07 - 0
    9
    T11: 14:45.582472132
    Z
    MSDN
    Magazine, 542
    2024 - 07 - 0
    9
    T11: 14:45.582473257
    Z
    mshtml.tlb(Microsoft
    HTML
    Object
    Library),
    2024 - 07 - 0
    9
    T11: 14:45.582474174
    Z
    835
    2024 - 07 - 0
    9
    T11: 14:45.582475049
    Z
    MSI(Microsoft
    Installer), 1079—1080
    2024 - 07 - 0
    9
    T11: 14:45.582475966
    Z
    MultiDataTrigger

    class , 321

    2024 - 07 - 0
    9
    T11: 14:45.582476841
    Z
    multiple
    document
    interface(MDI), 761, 1021—
    2024 - 07 - 0
    9
    T11: 14:45.582477799
    Z
    1022
    2024 - 07 - 0
    9
    T11: 14:45.582478674
    Z
    Multiselect
    property, 763
    2024 - 07 - 0
    9
    T11: 14:45.582479549
    Z
    multitargeting, 19
    2024 - 07 - 0
    9
    T11: 14:45.582480424
    Z
    multithreading
    2024 - 07 - 0
    9
    T11: 14:45.582481257
    Z
    BackgroundWorker
    component, 1045
    2024 - 07 - 0
    9
    T11: 14:45.582483966
    Z
    BeginInvoke()
    method, 1043—1044
    2024 - 07 - 0
    9
    T11: 14:45.582484882
    Z
    context, 1041
    2024 - 07 - 0
    9
    T11: 14:45.582485757
    Z
    CurrentDispatcher
    property, 1042
    2024 - 07 - 0
    9
    T11: 14:45.582486591
    Z
    definition
    of, 1041
    2024 - 07 - 0
    9
    T11: 14:45.582487507
    Z
    dispatcher, 1042
    2024 - 07 - 0
    9
    T11: 14:45.582488591
    Z
    DispatcherObject

    class , 1042

    2024 - 07 - 0
    9
    T11: 14:45.582489549
    Z
    DispatcherOperation
    object, 1044
    2024 - 07 - 0
    9
    T11: 14:45.582490424
    Z
    DispatcherPriority, 1044
    2024 - 07 - 0
    9
    T11: 14:45.582491382
    Z
    dual - core
    CPUs, 1041
    2024 - 07 - 0
    9
    T11: 14:45.582492257
    Z
    Invoke()
    method, 1045
    2024 - 07 - 0
    9
    T11: 14:45.582493132
    Z
    performing
    asynchronous
    operations, 1045
    2024 - 07 - 0
    9
    T11: 14:45.582494049
    Z
    performing
    time - consuming
    background
    2024 - 07 - 0
    9
    T11: 14:45.582495007
    Z
    operation, 1044
    2024 - 07 - 0
    9
    T11: 14:45.582495882
    Z
    single - threaded
    apartment
    model, 1041
    2024 - 07 - 0
    9
    T11: 14:45.582496841
    Z
    System.Threading.Thread
    object, 1045
    2024 - 07 - 0
    9
    T11: 14:45.582497882
    Z
    thread
    affinity, 1042
    2024 - 07 - 0
    9
    T11: 14:45.582498841
    Z
    thread
    rental, 1041
    2024 - 07 - 0
    9
    T11: 14:45.582499674
    Z
    VerifyAccess()
    method, 1043
    2024 - 07 - 0
    9
    T11: 14:45.582500549
    Z
    writing
    good
    multithreading
    code, 1045
    2024 - 07 - 0
    9
    T11: 14:45.582501466
    Z
    multitouch
    input, levels
    of
    support
    for , 149—157
    2024 - 07 - 0
    9
    T11: 14:45.582502549
    Z
    MultiTrigger

    class , 321, 323

    2024 - 07 - 0
    9
    T11: 14:45.582503424
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582504257
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582505132
    Z - ----
    2024 - 07 - 0
    9
    T11: 14:45.582506049
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582506882
    Z
    MustInherit

    class , 1057, 1064

    2024 - 07 - 0
    9
    T11: 14:45.582507757
    Z
    mutex, definition
    of, 228
    2024 - 07 - 0
    9
    T11: 14:45.582508632
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582509507
    Z ■ ** N **
    2024 - 07 - 0
    9
    T11: 14:45.582510382
    Z
    Name
    attribute, 31
    2024 - 07 - 0
    9
    T11: 14:45.582511216
    Z
    Name
    property, 269, 554
    2024 - 07 - 0
    9
    T11: 14:45.582512091
    Z
    namespaces
    2024 - 07 - 0
    9
    T11: 14:45.582513007
    Z
    core
    WPF
    namespace, 29
    2024 - 07 - 0
    9
    T11: 14:45.582516382
    Z
    core
    XAML
    namespace, 29
    2024 - 07 - 0
    9
    T11: 14:45.582517507
    Z
    declaring in XML, 29
    2024 - 07 - 0
    9
    T11: 14:45.582518382
    Z
    defining in XAML, 28
    2024 - 07 - 0
    9
    T11: 14:45.582519382
    Z.NET and, 29
    2024 - 07 - 0
    9
    T11: 14:45.582520257
    Z
    System.Windows.Shapes, 174
    2024 - 07 - 0
    9
    T11: 14:45.582521132
    Z
    using
    namespace
    prefixes, 47
    2024 - 07 - 0
    9
    T11: 14:45.582522007
    Z in WPF, 29
    2024 - 07 - 0
    9
    T11: 14:45.582522882
    Z
    XML
    namespaces as URIs, 29
    2024 - 07 - 0
    9
    T11: 14:45.582523757
    Z
    Narrator
    screen
    reader, 885
    2024 - 07 - 0
    9
    T11: 14:45.582524716
    Z
    native
    resolution, 8
    2024 - 07 - 0
    9
    T11: 14:45.582525591
    Z
    Navigate()
    method, 807—808, 834
    2024 - 07 - 0
    9
    T11: 14:45.582528132
    Z
    Navigated
    event, NavigationService

    class , 809

    2024 - 07 - 0
    9
    T11: 14:45.582533466
    Z
    NavigateToStream()
    method, WebBrowser
    2024 - 07 - 0
    9
    T11: 14:45.582534799
    Z

    class , 834

    2024 - 07 - 0
    9
    T11: 14:45.582535674
    Z
    NavigateToString()
    method, WebBrowser

    class ,
        2024 - 07 - 0
        9
        T11: 14:45.582536591
        Z
        834, 840

    2024 - 07 - 0
    9
    T11: 14:45.582537466
    Z
    NavigateUri
    property, 796, 950
    2024 - 07 - 0
    9
    T11: 14:45.582538382
    Z
    Navigating
    event, NavigationService

    class , 809

    2024 - 07 - 0
    9
    T11: 14:45.582539257
    Z
    NavigatingProgress, NavigationService

    class ,
        2024 - 07 - 0
        9
        T11: 14:45.582540216
        Z
        809

    2024 - 07 - 0
    9
    T11: 14:45.582541341
    Z
    NavigationCommands

    class , types of included

    2024 - 07 - 0
    9
    T11: 14:45.582542299
    Z
    commands, 270
    2024 - 07 - 0
    9
    T11: 14:45.582543174
    Z
    NavigationFailed
    event, 797, 809
    2024 - 07 - 0
    9
    T11: 14:45.582544091
    Z
    NavigationService

    class
        2024 - 07 - 0
        9
        T11: 14:45.582544966
        Z
        AddBackEntry()
        method, 810—811

    2024 - 07 - 0
    9
    T11: 14:45.582545966
    Z
    AddBackReference()
    method, 813—815
    2024 - 07 - 0
    9
    T11: 14:45.582546841
    Z
    adding
    custom
    items
    to
    journal, 811
    2024 - 07 - 0
    9
    T11: 14:45.582547799
    Z
    Application

    class , 808

    2024 - 07 - 0
    9
    T11: 14:45.582548674
    Z
    building
    linear
    navigation - based
    2024 - 07 - 0
    9
    T11: 14:45.582549632
    Z
    application, 809
    2024 - 07 - 0
    9
    T11: 14:45.582550507
    Z
    CanGoBack
    property, 807, 810
    2024 - 07 - 0
    9
    T11: 14:45.582551382
    Z
    CanGoForward
    property, 807
    2024 - 07 - 0
    9
    T11: 14:45.582552257
    Z
    Content
    property, 812
    2024 - 07 - 0
    9
    T11: 14:45.582553174
    Z
    creating
    page
    object
    manually, 807
    2024 - 07 - 0
    9
    T11: 14:45.582554049
    Z
    events
    of, 809
    2024 - 07 - 0
    9
    T11: 14:45.582555007
    Z
    ExtraData
    property, 808
    2024 - 07 - 0
    9
    T11: 14:45.582555882
    Z
    GetContentState()
    method, 813, 815
    2024 - 07 - 0
    9
    T11: 14:45.582556757
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582557591
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582558466
    Z
    GoBack()
    method, 807
    2024 - 07 - 0
    9
    T11: 14:45.582559341
    Z
    GoForward()
    method, 807
    2024 - 07 - 0
    9
    T11: 14:45.582560174
    Z
    Handled
    property, 808
    2024 - 07 - 0
    9
    T11: 14:45.582561091
    Z
    how
    WPF
    navigation
    occurs, 808
    2024 - 07 - 0
    9
    T11: 14:45.582562007
    Z
    InitializeComponent()
    method, 807
    2024 - 07 - 0
    9
    T11: 14:45.582562924
    Z
    IProvideCustomContentState
    interface, 813—
    2024 - 07 - 0
    9
    T11: 14:45.582563841
    Z
    814
    2024 - 07 - 0
    9
    T11: 14:45.582564716
    Z
    JournalEntryName
    property, 812
    2024 - 07 - 0
    9
    T11: 14:45.582565632
    Z
    ListSelectionJournalEntry
    callback, 813
    2024 - 07 - 0
    9
    T11: 14:45.582566549
    Z
    methods
    for controlling navigation stack,
    2024 - 07 - 0
    9
    T11: 14:45.582567466
    Z
    810
    2024 - 07 - 0
    9
    T11: 14:45.582568299
    Z
    Navigate()
    method, 807—808
    2024 - 07 - 0
    9
    T11: 14:45.582569299
    Z
    navigating
    to
    page
    based
    on
    its
    URI, 807
    2024 - 07 - 0
    9
    T11: 14:45.582570341
    Z
    RemoveBackEntry()
    method, 810
    2024 - 07 - 0
    9
    T11: 14:45.582572757
    Z
    Replay()
    method, 812, 814
    2024 - 07 - 0
    9
    T11: 14:45.582573632
    Z
    ReplayListChange
    delegate, 813
    2024 - 07 - 0
    9
    T11: 14:45.582574549
    Z
    returning
    information
    from page,
    816
    2024 - 07 - 0
    9
    T11: 14:45.582575424
    Z
    SourceItems
    property, 813
    2024 - 07 - 0
    9
    T11: 14:45.582576341
    Z
    StopLoading()
    method, 807
    2024 - 07 - 0
    9
    T11: 14:45.582577216
    Z
    suppressing
    navigation
    events, 808
    2024 - 07 - 0
    9
    T11: 14:45.582578132
    Z
    table
    of
    navigation
    events, 808
    2024 - 07 - 0
    9
    T11: 14:45.582579007
    Z
    TargetItems
    property, 813
    2024 - 07 - 0
    9
    T11: 14:45.582579882
    Z
    using
    Refresh()
    to
    reload
    page, 807
    2024 - 07 - 0
    9
    T11: 14:45.582580757
    Z
    WPF
    navigation as asynchronous, 807
    2024 - 07 - 0
    9
    T11: 14:45.582581674
    Z
    NavigationService
    property, Page

    class , 795

    2024 - 07 - 0
    9
    T11: 14:45.582582549
    Z
    NavigationStopped
    event, NavigationService
    2024 - 07 - 0
    9
    T11: 14:45.582583466
    Z

    class , 809

    2024 - 07 - 0
    9
    T11: 14:45.582584341
    Z
    NavigationUIVisibility
    property, 800—801
    2024 - 07 - 0
    9
    T11: 14:45.582585299
    Z
    NavigationWindow

    class , 793, 797

    2024 - 07 - 0
    9
    T11: 14:45.582586424
    Z
    NearPlaneDistance
    property, 903
    2024 - 07 - 0
    9
    T11: 14:45.582587382
    Z.NET
    2024 - 07 - 0
    9
    T11: 14:45.582588216
    Z
    Code
    DOM
    model, 54
    2024 - 07 - 0
    9
    T11: 14:45.582589132
    Z
    global assembly
    cache(GAC), 242
    2024 - 07 - 0
    9
    T11: 14:45.582590007
    Z
    ildasm, 235
    2024 - 07 - 0
    9
    T11: 14:45.582590882
    Z
    mapping.NET
    namespace
    to
    XML
    2024 - 07 - 0
    9
    T11: 14:45.582591799
    Z
    namespace, 46
    2024 - 07 - 0
    9
    T11: 14:45.582592799
    Z
    namespaces in, 29
    2024 - 07 - 0
    9
    T11: 14:45.582593674
    Z
    p / invoke, 769
    2024 - 07 - 0
    9
    T11: 14:45.582594549
    Z
    probing, 240
    2024 - 07 - 0
    9
    T11: 14:45.582595382
    Z
    replacing.NET
    properties
    with dependency
        2024 - 07 - 0
        9
        T11: 14:45.582596257
        Z
        properties, 105
    2024 - 07 - 0
    9
    T11: 14:45.582597174
    Z
    ResourceManager

    class , 236

    2024 - 07 - 0
    9
    T11: 14:45.582598049
    Z
    ResourceSet

    class , 236

    2024 - 07 - 0
    9
    T11: 14:45.582598924
    Z
    satellite
    assemblies, 240
    2024 - 07 - 0
    9
    T11: 14:45.582599799
    Z
    type converters, 34
    2024 - 07 - 0
    9
    T11: 14:45.582600674
    Z
    window
    ownership, 760
    2024 - 07 - 0
    9
    T11: 14:45.582601549
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582602341
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582603174
    Z - ----
    2024 - 07 - 0
    9
    T11: 14:45.582604049
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582604882
    Z
    XML
    capabilities in, 640
    2024 - 07 - 0
    9
    T11: 14:45.582605716
    Z.NET
    1.
    x, 61
    2024 - 07 - 0
    9
    T11: 14:45.582606632
    Z.NET
    2.0, 1
    2024 - 07 - 0
    9
    T11: 14:45.582607466
    Z
    BackgroundWorker
    component, 1045
    2024 - 07 - 0
    9
    T11: 14:45.582608466
    Z
    coordinate - based
    layout, 62
    2024 - 07 - 0
    9
    T11: 14:45.582609299
    Z
    enhancing
    Button and Label
    classes, 175
    2024 - 07 - 0
    9
    T11: 14:45.582610216
    Z
    flow - based
    layout
    panels, 62
    2024 - 07 - 0
    9
    T11: 14:45.582611382
    Z
    FlowLayoutPanel, 61
    2024 - 07 - 0
    9
    T11: 14:45.582613507
    Z
    SoundPlayer

    class , 865

    2024 - 07 - 0
    9
    T11: 14:45.582614382
    Z
    System.Drawing
    namespace, 301
    2024 - 07 - 0
    9
    T11: 14:45.582615299
    Z
    System.Media.SystemSounds

    class , 868

    2024 - 07 - 0
    9
    T11: 14:45.582616174
    Z
    TableLayoutPanel, 61
    2024 - 07 - 0
    9
    T11: 14:45.582617049
    Z.NET
    Framework
    3.0, 17
    2024 - 07 - 0
    9
    T11: 14:45.582617966
    Z
    no - argument
    constructors, 47
    2024 - 07 - 0
    9
    T11: 14:45.582618841
    Z
    nonclient
    area, definition
    of, 751
    2024 - 07 - 0
    9
    T11: 14:45.582619716
    Z
    nonrectangular
    windows
    2024 - 07 - 0
    9
    T11: 14:45.582620591
    Z
    adding
    sizing
    grip
    to
    shaped
    window, 769
    2024 - 07 - 0
    9
    T11: 14:45.582621507
    Z
    comparing
    background - based and shapedrawing
    approaches, 767
    2024 - 07 - 0
    9
    T11: 14:45.582622507
    Z
    creating
    shaped
    window
    with rounded
        2024 - 07 - 0
        9
        T11: 14:45.582623424
        Z
        Border
        element, 765
    2024 - 07 - 0
    9
    T11: 14:45.582624341
    Z
    creating
    transparent
    window, 764, 766
    2024 - 07 - 0
    9
    T11: 14:45.582632549
    Z
    detecting
    mouse
    movements
    over
    edges
    of
    2024 - 07 - 0
    9
    T11: 14:45.582633466
    Z
    window, 769
    2024 - 07 - 0
    9
    T11: 14:45.582634341
    Z
    initiating
    window
    dragging
    mode
    by
    calling
    2024 - 07 - 0
    9
    T11: 14:45.582637882
    Z
    Window.DragMove(), 768
    2024 - 07 - 0
    9
    T11: 14:45.582638924
    Z
    moving
    shaped
    windows, 768
    2024 - 07 - 0
    9
    T11: 14:45.582639841
    Z
    placing
    Rectangle
    that
    allows
    right - side
    2024 - 07 - 0
    9
    T11: 14:45.582641007
    Z
    window
    resizing, 769
    2024 - 07 - 0
    9
    T11: 14:45.582641882
    Z
    placing
    sizing
    grip
    correctly, 769
    2024 - 07 - 0
    9
    T11: 14:45.582642799
    Z
    procedure
    for creating shaped window, 763
    2024 - 07 - 0
    9
    T11: 14:45.582643716
    Z
    providing
    background
    art, 764
    2024 - 07 - 0
    9
    T11: 14:45.582644591
    Z
    removing
    standard
    window
    appearance
    2024 - 07 - 0
    9
    T11: 14:45.582645466
    Z(window
    chrome), 764
    2024 - 07 - 0
    9
    T11: 14:45.582646341
    Z
    resizing
    shaped
    windows, 769
    2024 - 07 - 0
    9
    T11: 14:45.582647299
    Z
    resizing
    window
    manually
    by
    setting
    its
    2024 - 07 - 0
    9
    T11: 14:45.582648174
    Z
    Width
    property, 769
    2024 - 07 - 0
    9
    T11: 14:45.582649049
    Z
    setting
    Window.ResizeMode
    property, 769
    2024 - 07 - 0
    9
    T11: 14:45.582650007
    Z
    using
    Path
    element
    to
    create
    background,
    2024 - 07 - 0
    9
    T11: 14:45.582650882
    Z
    767
    2024 - 07 - 0
    9
    T11: 14:45.582651757
    Z
    Nonzero
    fill
    rule, 347
    2024 - 07 - 0
    9
    T11: 14:45.582652632
    Z
    normal
    2024 - 07 - 0
    9
    T11: 14:45.582653466
    Z
    calculating
    normal
    that
    's perpendicular to
    2024 - 07 - 0
    9
    T11: 14:45.582654424
    Z
    triangle
    's surface, 908
    2024 - 07 - 0
    9
    T11: 14:45.582655299
    Z
    definition
    of, 906
    2024 - 07 - 0
    9
    T11: 14:45.582656132
    Z
    guidelines
    for choosing right normals, 908
    2024 - 07 - 0
    9
    T11: 14:45.582657049
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582657924
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582658757
    Z
    problem
    of
    sharing
    Position
    points and
    2024 - 07 - 0
    9
    T11: 14:45.582659674
    Z
    sharing
    normals, 907
    2024 - 07 - 0
    9
    T11: 14:45.582662299
    Z
    understanding, 906
    2024 - 07 - 0
    9
    T11: 14:45.582663424
    Z
    Normal
    state, controls, 541, 572
    2024 - 07 - 0
    9
    T11: 14:45.582664341
    Z
    Normals
    property, 893—894, 909
    2024 - 07 - 0
    9
    T11: 14:45.582665382
    Z
    NotifyIcon

    class , 1025

    2024 - 07 - 0
    9
    T11: 14:45.582666257
    Z
    NotifyOnValidationError
    property, 623, 628
    2024 - 07 - 0
    9
    T11: 14:45.582667132
    Z
    null
    markup
    extension, 179
    2024 - 07 - 0
    9
    T11: 14:45.582668049
    Z
    NullExtension, 38
    2024 - 07 - 0
    9
    T11: 14:45.582668924
    Z
    NumericUpDown
    control, 1020
    2024 - 07 - 0
    9
    T11: 14:45.582669841
    Z
    2024 - 07 - 0
    9
    T11: 14:45.582670757
    Z ■ ** O **
    2024 - 07 - 0
    9
    T11: 14:45.582671632
    Z
    object
    resources
    2024 - 07 - 0
    9
    T11: 14:45.582672549
    Z
    accessing
    resources in code, 299
    2024 - 07 - 0
    9
    T11: 14:45.582673424
    Z
    adding
    resources
    programmatically, 300
    2024 - 07 - 0
    9
    T11: 14:45.582674299
    Z
    advantages
    of, 293
    2024 - 07 - 0
    9
    T11: 14:45.582675132
    Z
    application
    resources, 300
    2024 - 07 - 0
    9
    T11: 14:45.582676007
    Z
    ComponentResourceKey, 305
    2024 - 07 - 0
    9
    T11: 14:45.582676924
    Z
    creating
    resource
    dictionary, 302
    2024 - 07 - 0
    9
    T11: 14:45.582677841
    Z
    defining
    image
    brush as resource, 294
    2024 - 07 - 0
    9
    T11: 14:45.582678757
    Z
    defining
    resources
    at
    window
    level, 294
    2024 - 07 - 0
    9
    T11: 14:45.582679632
    Z
    definition
    of, 293
    2024 - 07 - 0
    9
    T11: 14:45.582680466
    Z
    FrameworkElement.FindResource()
    2024 - 07 - 0
    9
    T11: 14:45.582681341
    Z
    method, 299
    2024 - 07 - 0
    9
    T11: 14:45.582682174
    Z
    Freezable

    class , 297

    2024 - 07 - 0
    9
    T11: 14:45.582683049
    Z
    generic.xaml
    file, code
    example, 306
    2024 - 07 - 0
    9
    T11: 14:45.582683924
    Z
    hierarchy
    of
    resources, 295
    2024 - 07 - 0
    9
    T11: 14:45.582684841
    Z
    ImageSource
    property, 306
    2024 - 07 - 0
    9
    T11: 14:45.582685966
    Z
    Key
    attribute, 294
    2024 - 07 - 0
    9
    T11: 14:45.582686882
    Z
    nonshared
    resources, reasons
    for using, 299
        2024 - 07 - 0
        9
        T11: 14:45.582688216
        Z
        resource
        keys, 301
    2024 - 07 - 0
    9
    T11: 14:45.582689091
    Z
    ResourceDictionary

    class , 294

    2024 - 07 - 0
    9
    T11: 14:45.582689966
    Z
    ResourceKey
    property, 296
    2024 - 07 - 0
    9
    T11: 14:45.582690882
    Z
    resources
    collection, 294
    2024 - 07 - 0
    9
    T11: 14:45.582691799
    Z
    Resources
    property, 294
    2024 - 07 - 0
    9
    T11: 14:45.582692757
    Z
    reusing
    resource
    names, 296
    2024 - 07 - 0
    9
    T11: 14:45.582693632
    Z
    sharing
    resources
    among
    assemblies, 304
    2024 - 07 - 0
    9
    T11: 14:45.582694591
    Z
    static
    vs.
    2024 - 07 - 0
    9
    T11: 14:45.582695466
    Z'''

    chunk_list = [ChunkedItem(content=large_content, url="http://example.com", document_id=1, parent_id=1)]

    # Define collection name, model name, and model dimension
    collection_name = "user_625ece7e042d4f40bd2588b16bec7be6"
    model_name = "paraphrase-multilingual-mpnet-base-v2"
    model_dimension = 768

    # Call the store_chunk_list method
    try:
        milvus_db.store_chunk_list(chunk_list, collection_name, model_name, model_dimension)
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    test_store_chunk_list()