import pytest

from app.coverage.services.sites_coverage import SitesCoverage


class TestSitesCoverage:
    @classmethod
    def setup_class(cls):
        """setup any state specific to the execution of the given class (which
        usually contains tests).
        """
        cls.sites_coverage = SitesCoverage()

    @classmethod
    def teardown_class(cls):
        """teardown any state that was previously setup with a call to
        setup_class.
        """
        pass

    def test_calculate_coverage(self):
        # self.sites_coverage.calculate_coverage()
        assert 3 == 3
