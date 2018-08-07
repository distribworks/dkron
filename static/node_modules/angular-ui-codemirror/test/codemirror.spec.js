describe('uiCodemirror', function() {
  'use strict';

  // declare these up here to be global to all tests
  var scope, $compile, $timeout, uiConfig;
  var codemirrorDefaults = window.CodeMirror.defaults;

  beforeEach(function() {
    module('ui.codemirror');

    inject(function(_$rootScope_, _$compile_, _$timeout_, uiCodemirrorConfig) {
      scope = _$rootScope_.$new();
      $compile = _$compile_;
      $timeout = _$timeout_;
      uiConfig = uiCodemirrorConfig;
    });

  });

  afterEach(function() {
    uiConfig = {};
  });


  it('should not throw an error when window.CodeMirror is defined', function() {
    function compile() {
      $compile('<div ui-codemirror></div>')(scope);
    }

    var _CodeMirror = window.CodeMirror;
    delete window.CodeMirror;
    expect(window.CodeMirror).toBeUndefined();
    expect(compile)
      .toThrow(new Error('ui-codemirror needs CodeMirror to work... (o rly?)'));
    window.CodeMirror = _CodeMirror;
  });

  describe('destruction', function() {

    var parentElement;

    beforeEach(function() {
      parentElement = angular.element('<div></div>');
      angular.element(document.body).prepend(parentElement);
    });

    afterEach(function() {
      parentElement.remove();
    });

    function shouldDestroyTest(elementType, template) {
      it('should destroy the directive of ' + elementType, function() {
        var element = angular.element(template);
        parentElement.append(element);

        $compile(element)(scope);
        scope.$digest();

        expect(parentElement.children().length).toBe(1);
        element.remove();
        scope.$digest();
        expect(parentElement.children().length).toBe(0);
      });
    }

    shouldDestroyTest('an element', '<ui-codemirror></ui-codemirror>');
    shouldDestroyTest('an attribute', '<div ui-codemirror=""></div>');

  });

  it('should not throw an error when window.CodeMirror is defined an attribute', function() {
    function compile() {
      $compile('<div ui-codemirror></div>')(scope);
    }

    expect(window.CodeMirror).toBeDefined();
    expect(compile).not.toThrow();
  });


  it('should not throw an error when window.CodeMirror is defined an element', function() {
    function compile() {
      $compile('<ui-codemirror></ui-codemirror>')(scope);
    }

    expect(window.CodeMirror).toBeDefined();
    expect(compile).not.toThrow();
  });


  it('should watch all uiCodemirror attribute', function() {
    spyOn(scope, '$watch');
    scope.cmOption = {};
    $compile('<div ui-codemirror="cmOption"  ng-model="foo" ui-refresh="sdf"></div>')(scope);
    expect(scope.$watch.calls.count()).toEqual(3); // The uiCodemirror+ the ngModel + the uiRefresh
    expect(scope.$watch).toHaveBeenCalledWith('cmOption', jasmine.any(Function), true); // uiCodemirror
    expect(scope.$watch).toHaveBeenCalledWith(jasmine.any(Function)); // ngModel
    expect(scope.$watch).toHaveBeenCalledWith('sdf', jasmine.any(Function)); // uiRefresh
  });

  describe('CodeMirror instance', function() {

    var codemirror = null, spies = angular.noop;

    beforeEach(function() {
      var _constructor = window.CodeMirror;
      window.CodeMirror = jasmine.createSpy('window.CodeMirror')
        .and.callFake(function() {
          codemirror = _constructor.apply(this, arguments);
          spies(codemirror);
          return codemirror;
        });

      window.CodeMirror.defaults = codemirrorDefaults;
    });


    it('should call the CodeMirror constructor with a function', function() {
      $compile('<div ui-codemirror></div>')(scope);

      expect(window.CodeMirror.calls.count()).toEqual(1);
      expect(window.CodeMirror)
        .toHaveBeenCalledWith(jasmine.any(Function), jasmine.any(Object));

      expect(codemirror).toBeDefined();
    });

    it('should work as an element', function() {
      $compile('<ui-codemirror></ui-codemirror>')(scope);

      expect(window.CodeMirror.calls.count()).toEqual(1);
      expect(window.CodeMirror)
        .toHaveBeenCalledWith(jasmine.any(Function), jasmine.any(Object));

      expect(codemirror).toBeDefined();
    });

    it('should have a child element with a div.CodeMirror', function() {
      // Explicit a parent node to support the directive.
      var element = $compile('<div ui-codemirror></div>')(scope).children();

      expect(element).toBeDefined();
      expect(element.prop('tagName')).toBe('DIV');
      expect(element.prop('classList').length).toEqual(2);
      expect(element.prop('classList')[0]).toEqual('CodeMirror');
      expect(element.prop('classList')[1]).toEqual('cm-s-default');
    });


    describe('options', function() {

      spies = function(codemirror) {
        codemirror._setOption = codemirror._setOption || codemirror.setOption;
        codemirror.setOption = jasmine.createSpy('codemirror.setOption')
          .and.callFake(function() {
            codemirror._setOption.apply(this, arguments);
          });
      };

      it('should not be called', function() {
        $compile('<div ui-codemirror></div>')(scope);
        expect(window.CodeMirror)
          .toHaveBeenCalledWith(jasmine.any(Function), { value: '' });
        expect(codemirror.setOption).not.toHaveBeenCalled();
      });

      it('should include the passed options (attribute directive)', function() {
        $compile('<div ui-codemirror="{oof: \'baar\'}"></div>')(scope);

        expect(window.CodeMirror)
          .toHaveBeenCalledWith(jasmine.any(Function), {
            value: '',
            oof: 'baar'
          });
        expect(codemirror.setOption).not.toHaveBeenCalled();
      });

      it('should include the passed options (element directive)', function() {
        $compile('<ui-codemirror ui-codemirror-opts="{oof: \'baar\'}"></ui-codemirror>')(scope);

        expect(window.CodeMirror)
          .toHaveBeenCalledWith(jasmine.any(Function), {
            value: '',
            oof: 'baar'
          });
        expect(codemirror.setOption).not.toHaveBeenCalled();
      });

      it('should include the default options', function() {
        uiConfig.codemirror = { bar: 'baz' };
        $compile('<div ui-codemirror></div>')(scope);

        expect(window.CodeMirror).toHaveBeenCalledWith(jasmine.any(Function), {
          value: '',
          bar: 'baz'
        });
        expect(codemirror.setOption).not.toHaveBeenCalled();
      });

      it('should extent the default options', function() {
        uiConfig.codemirror = { bar: 'baz' };
        $compile('<div ui-codemirror="{oof: \'baar\'}"></div>')(scope);

        expect(window.CodeMirror).toHaveBeenCalledWith(jasmine.any(Function), {
          value: '',
          oof: 'baar',
          bar: 'baz'
        });
        expect(codemirror.setOption).not.toHaveBeenCalled();
      });

      it('should impact codemirror', function() {
        uiConfig.codemirror = {};
        $compile('<div ui-codemirror="{theme: \'baar\'}"></div>')(scope);

        expect(window.CodeMirror).toHaveBeenCalledWith(jasmine.any(Function), {
          value: '',
          theme: 'baar'
        });
        expect(codemirror.setOption).not.toHaveBeenCalled();

        expect(codemirror.getOption('theme')).toEqual('baar');
      });
    });

    it('should not trigger watch ui-refresh', function() {
      spyOn(scope, '$watch');
      $compile('<div ui-codemirror ui-refresh=""></div>')(scope);
      expect(scope.$watch).not.toHaveBeenCalled();
    });

    it('should trigger the CodeMirror.refresh() method', function() {
      $compile('<div ui-codemirror ui-refresh="bar"></div>')(scope);


      spyOn(codemirror, 'refresh');
      scope.$apply('bar = null');

      scope.$apply('bar = false');
      expect(scope.bar).toBeFalsy();
      $timeout.flush();
      expect(codemirror.refresh).toHaveBeenCalled();
      scope.$apply('bar = true');
      expect(scope.bar).toBeTruthy();
      $timeout.flush();
      expect(codemirror.refresh).toHaveBeenCalled();
      scope.$apply('bar = 0');
      expect(scope.bar).toBeFalsy();
      $timeout.flush();
      expect(codemirror.refresh).toHaveBeenCalled();
      scope.$apply('bar = 1');
      expect(scope.bar).toBeTruthy();
      $timeout.flush();
      expect(codemirror.refresh).toHaveBeenCalled();

      expect(codemirror.refresh.calls.count()).toEqual(4);
    });


    it('when the IDE changes should update the model', function() {
      var element = $compile('<div ui-codemirror ng-model="foo"></div>')(scope);
      var ctrl = element.controller('ngModel');

      expect(ctrl.$pristine).toBe(true);
      expect(ctrl.$valid).toBe(true);

      var value = 'baz';
      codemirror.setValue(value);
      scope.$apply();
      expect(scope.foo).toBe(value);

      expect(ctrl.$valid).toBe(true);
      expect(ctrl.$dirty).toBe(true);

    });

    it('when the model changes should update the IDE', function() {
      var element = $compile('<div ui-codemirror ng-model="foo"></div>')(scope);
      var ctrl = element.controller('ngModel');

      expect(ctrl.$pristine).toBe(true);
      expect(ctrl.$valid).toBe(true);

      scope.$apply('foo = "bar"');
      expect(codemirror.getValue()).toBe(scope.foo);

      expect(ctrl.$pristine).toBe(true);
      expect(ctrl.$valid).toBe(true);
    });


    it('when the IDE changes should use ngChange', function() {
      scope.change = angular.noop;
      spyOn(scope, 'change').and.callFake(function() { expect(scope.foo).toBe('baz'); });

      $compile('<div ui-codemirror ng-model="foo" ng-change="change()"></div>')(scope);

      // change shouldn't be called initialy
      expect(scope.change).not.toHaveBeenCalled();


      // change shouldn't be called when the value change is coming from the model.
      scope.$apply('foo = "bar"');
      expect(scope.change).not.toHaveBeenCalled();

      // change should be called when user changes the input.
      codemirror.setValue('baz');
      scope.$apply();
      expect(scope.change.calls.count()).toBe(1);
      expect(scope.change).toHaveBeenCalledWith();
    });

    it('should runs the onLoad callback', function() {
      scope.codemirrorLoaded = jasmine.createSpy('scope.codemirrorLoaded');

      $compile('<div ui-codemirror="{onLoad: codemirrorLoaded}"></div>')(scope);

      expect(scope.codemirrorLoaded).toHaveBeenCalled();
      expect(scope.codemirrorLoaded).toHaveBeenCalledWith(codemirror);
    });

    it('responds to the $broadcast event "CodeMirror"', function() {
      var broadcast = {};
      broadcast.callback = jasmine.createSpy('broadcast.callback');

      $compile('<div ui-codemirror></div>')(scope);
      scope.$broadcast('CodeMirror', broadcast.callback);

      expect(broadcast.callback).toHaveBeenCalled();
      expect(broadcast.callback).toHaveBeenCalledWith(codemirror);
    });


    it('should watch the options (attribute directive)', function() {

      scope.cmOption = { readOnly: true };
      $compile('<div ui-codemirror="cmOption"></div>')(scope);
      scope.$digest();

      expect(codemirror.getOption('readOnly')).toBeTruthy();

      scope.cmOption.readOnly = false;
      scope.$digest();
      expect(codemirror.getOption('readOnly')).toBeFalsy();
    });

    it('should watch the options (element directive)', function() {

      scope.cmOption = { readOnly: true };
      $compile('<ui-codemirror ui-codemirror-opts="cmOption"></div>')(scope);
      scope.$digest();

      expect(codemirror.getOption('readOnly')).toBeTruthy();

      scope.cmOption.readOnly = false;
      scope.$digest();
      expect(codemirror.getOption('readOnly')).toBeFalsy();
    });

    it('should watch the options (object property)', function() {

      scope.cm = {};
      scope.cm.option = { readOnly: true };
      $compile('<div ui-codemirror="cm.option"></div>')(scope);
      scope.$digest();

      expect(codemirror.getOption('readOnly')).toBeTruthy();

      scope.cm.option.readOnly = false;
      scope.$digest();
      expect(codemirror.getOption('readOnly')).toBeFalsy();
    });

  });

  it('when the model is an object or an array should throw an error', function() {
    function compileWithObject() {
      $compile('<div ui-codemirror ng-model="foo"></div>')(scope);
      scope.foo = {};
      scope.$apply();
    }

    function compileWithArray() {
      $compile('<div ui-codemirror ng-model="foo"></div>')(scope);
      scope.foo = [];
      scope.$apply();
    }

    expect(compileWithObject).toThrow();
    expect(compileWithArray).toThrow();
  });

});
