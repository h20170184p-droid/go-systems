import 'package:flutter/material.dart';
import 'package:pdfx/pdfx.dart';
import 'package:file_picker/file_picker.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return const MaterialApp(
      debugShowCheckedModeBanner: false,
      home: PdfReaderPage(),
    );
  }
}

// --- Linked List ---
class PdfList {
  Map<String, List<int>> nodes = {};
  int count = 0;

  void add(String path) {
    if (nodes.containsKey(path)) return;
    if (nodes.isEmpty) {
      nodes[path] = [0, 1];
    } else {
      for (var val in nodes.values) {
        if (val[1] == count) {
          nodes[path] = [val[1], val[1] + 1];
          break;
        }
      }
    }
    count++;
  }

  void remove(String path) {
    if (!nodes.containsKey(path)) return;
    List<int> a = nodes[path]!;
    nodes.remove(path);
    for (var k in nodes.keys.toList()) {
      if (nodes[k]![0] >= a[0]) {
        nodes[k] = [nodes[k]![0] - 1, nodes[k]![1] - 1];
      }
    }
    count--;
  }

  List<String> ordered() {
    List<String> result = List.filled(nodes.length, '');
    for (var entry in nodes.entries) {
      int index = entry.value[1] - 1;
      if (index >= 0 && index < result.length) {
        result[index] = entry.key;
      }
    }
    return result;
  }
}

// --- Main Page ---
class PdfReaderPage extends StatefulWidget {
  const PdfReaderPage({super.key});

  @override
  State<PdfReaderPage> createState() => _PdfReaderPageState();
}

class _PdfReaderPageState extends State<PdfReaderPage> {
  final PdfList _pdfList = PdfList();
  int _currentIndex = 0;
  PdfControllerPinch? _controller;

  // Stores last visited page for each PDF path
  final Map<String, int> _savedPages = {};

  List<String> get _ordered => _pdfList.ordered();

  void _saveCurrentPage() {
    if (_controller == null) return;
    final ordered = _ordered;
    if (ordered.isEmpty) return;
    final path = ordered[_currentIndex];
    _savedPages[path] = _controller!.page;
  }

  void _loadPdf(String path) {
    _controller?.dispose();
    final savedPage = _savedPages[path] ?? 1;
    _controller = PdfControllerPinch(
      document: PdfDocument.openFile(path),
      initialPage: savedPage,
    );
  }

  void _unloadPdf() {
    _controller?.dispose();
    _controller = null;
  }

  Future<void> _addPdf() async {
    FilePickerResult? result = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowedExtensions: ['pdf'],
    );
    if (result == null) return;
    String path = result.files.single.path!;
    if (_pdfList.nodes.containsKey(path)) return;

    setState(() {
      _saveCurrentPage();
      _pdfList.add(path);
      _currentIndex = _pdfList.count - 1;
      _loadPdf(path);
    });
  }

  void _closeCurrent() {
    if (_ordered.isEmpty) return;
    String path = _ordered[_currentIndex];
    _savedPages.remove(path);
    _pdfList.remove(path);

    setState(() {
      if (_pdfList.count == 0) {
        _unloadPdf();
      } else {
        if (_currentIndex >= _pdfList.count) {
          _currentIndex = _pdfList.count - 1;
        }
        _loadPdf(_ordered[_currentIndex]);
      }
    });
  }

  void _goLeft() {
    if (_currentIndex > 0) {
      setState(() {
        _saveCurrentPage();
        _currentIndex--;
        _loadPdf(_ordered[_currentIndex]);
      });
    }
  }

  void _goRight() {
    if (_currentIndex < _pdfList.count - 1) {
      setState(() {
        _saveCurrentPage();
        _currentIndex++;
        _loadPdf(_ordered[_currentIndex]);
      });
    }
  }

  @override
  void dispose() {
    _controller?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    List<String> ordered = _ordered;

    return Scaffold(
      appBar: AppBar(
        title: Text(
          ordered.isEmpty
              ? 'PDF Reader'
              : '${_currentIndex + 1} / ${ordered.length}',
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.arrow_back_ios),
            onPressed: _goLeft,
          ),
          IconButton(
            icon: const Icon(Icons.arrow_forward_ios),
            onPressed: _goRight,
          ),
          IconButton(
            icon: const Icon(Icons.close),
            onPressed: _closeCurrent,
          ),
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: _addPdf,
          ),
        ],
      ),
      body: ordered.isEmpty
          ? const Center(child: Text('No PDFs loaded. Tap + to add.'))
          : PdfViewPinch(
              key: ValueKey(ordered[_currentIndex]),
              controller: _controller!,
            ),
    );
  }
}
